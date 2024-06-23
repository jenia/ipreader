# IpReader

## Architecture description

This program reads a file that contains one IP address per line and counts the number of unique IP.

This program has two packages:
- ipreader
- ipcounter

### ipreader package
This package defines a function `ReadFile` that reads a file into a buffer and sends all the IPs contained in that buffer as a slice of strings to `Counter.AddIpSlice(ips []string)`.

Note the signature of the function `func ReadFile(file io.Reader, counter Counter, buffer []byte) error`

It takes
1. file handle `io.Reader`
2. a `Counter` interface
```
type Counter interface {
	AddIpSlice(ips []string)
	Close()
}
```
3. a `[]byte` buffer

Using an interface, I decouple `ReadFile` function from whatever implements the `Counter` interface. In this case it's `ipcounter.IpCounter` that asynchronously counts the IPs in the slice `AddIpSlice(ips []string)`.
This decoupling allows for simpler testing by unit testing each feature separately and simpler benchmarking but allowing to benchmark each feature separately.
See the tests and benchmarks for details.

Furthermore, using a buffer, we can read the file 1MB at a time or 10MB at a time and bulk process the IPs that buffer. Furthermore, we control how much memory the program consumes i.e. whatever we pick as the size of the buffer.


### ipcounter package
This package defines a struct `IpCounter` that counts unique IPs sent to through the function `AddIpSlice`. Each item in the slice is one IP.
The struct `IpCounter` is meant to be instantiated once.
The function `Count` can run, and it meant to be run (see benchmarking), on multiple go-routines

See the main function or the bench marking tests in `ipcounter` package to see how to instantiate and use `IpCounter`

## BenchMarking results.

It seems that this program has better performance with larger slices. This was expected.

```
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

### Conclusion regarding benchmarking


|                                          | ns/item              |
| ---------------------------------------- | -------------------- |
| BenchmarkCount1GoRoutine100ItemSlice-8   | 459.2 /100 = 4.592   |
| BenchmarkCount1GoRoutine1000ItemSlice-8  | 567.7/1000 = 0.567.7 |
| BenchmarkCount10GoRoutine100ItemSlice-8  | 535.3/100  = 5.353   |
| BenchmarkCount10GoRoutine1000ItemSlice-8 | 588.2/1000 = 0.5882  |


This means that it is best to initiate the IpReader with a bigger buffer and send bigger slices to IpCounter
