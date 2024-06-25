# Welcome to batch IP Reader
An asynchronous text processor and batch file reader.

## Table of contents
  - [Architecture description](#Architecture-description)
    - [ipreader package](#ipreader-package)
	- [ipcounter package](#ipcounter-package)
  - [Profiling](#Profiling)
    - [CPU Benchmarking results](#CPU-Benchmarking-results)
	- [Conclusion regarding CPU benchmarking](#Conclusion-regarding-CPU-benchmarking)

## Architecture description

This program reads a file that contains one IP address per line and asynchronously counts the number of unique IPs contained in that file.

This program has two packages:
- ipreader
- ipcounter

This program is carefully designed to process large volume of data while using a minimum of resources:

- This program is profile using Golang [pprof tooling](profiling)
- ips are treated as numbers rather than strings to optimize memory and lookup performance
- ipcounter pre-allocate a 2^32 slice to store all IP V4 to avoid doing unnecessary malloc and data copying
- ipcounter and ipreader are async to maximize CPU utilization
- ipreader is carefully designed to read the file one buffer size at a time to minimize context switching

### ipreader package
This package defines a function `ReadFile` that reads a file into a buffer and sends all the IPs contained in that buffer to `Counter.AddIpSlice(ips []uint32)`.

Note the signature of the function `func ReadFile(file io.Reader, counter Counter, buffer []byte) error`

It takes
1. `io.Reader` file handle
2. `Counter` interface
```
type Counter interface {
	AddIpSlice(ips []uint32)
	Close()
}
```
3. a `[]byte` buffer

Using an interface, I decouple `ReadFile` function from whatever implements the `Counter` interface which is Golang best practices. In this case it's `ipcounter.IpCounter` that implements the `Counter` interface by asynchronously counting the IPs in the slice `AddIpSlice(ips []uint32)`.
This decoupling allows for higher modularity which has important implications:
1. Simpler code
2. Simpler testing. Each feature can be unit tested separately. Same thing for benchmarking. See the tests and benchmarks for details.

Furthermore, using a buffer + goroutines we can bulk read and async process the IPs.

### ipcounter package

This package defines a struct `IpCounter` that counts unique IPs sent to it using the function `AddIpSlice`. Each item in the slice is one IP.
The struct `IpCounter` is meant to be instantiated once.
The function `Count` can run, and it meant to be run (see benchmarking), on multiple go-routines

See the main function or the bench marking tests in `ipcounter` package to see how to instantiate and use `IpCounter`

## Profiling

This program is meticulously profile using Golang pprof tooling

### CPU Benchmarking results

It seems that this program has better performance with larger slices. This was expected.

```
cd ipcounter
go test -bench=. -benchmem -memprofile mem.prof -cpuprofile cpu.prof -benchtime=10s

goos: linux
goarch: amd64
pkg: github.com/Ecwid/new-job/ipcounter
cpu: Intel(R) Core(TM) i7-4790K CPU @ 4.00GHz
BenchmarkCount1GoRoutine100ItemSlice-8     	33546852	       459.2 ns/op	      99 B/op	       3 allocs/op
BenchmarkCount1GoRoutine1000ItemSlice-8    	32292217	       567.7 ns/op	     145 B/op	       3 allocs/op
BenchmarkCount10GoRoutine100ItemSlice-8    	28279573	       535.3 ns/op	     107 B/op	       3 allocs/op
BenchmarkCount10GoRoutine1000ItemSlice-8   	33324129	       588.2 ns/op	     143 B/op	       3 allocs/op
PASS
ok  	github.com/Ecwid/new-job/ipcounter	70.597s
```

### Conclusion regarding CPU benchmarking
|                                          | ns/item              |
| ---------------------------------------- | -------------------- |
| BenchmarkCount1GoRoutine100ItemSlice-8   | 459.2 /100 = 4.592   |
| BenchmarkCount1GoRoutine1000ItemSlice-8  | 567.7/1000 = 0.567.7 |
| BenchmarkCount10GoRoutine100ItemSlice-8  | 535.3/100  = 5.353   |
| BenchmarkCount10GoRoutine1000ItemSlice-8 | 588.2/1000 = 0.5882  |


This means that it is best to initiate the IpReader with a bigger buffer and send bigger slices to IpCounter

