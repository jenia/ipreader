# Welcome to batch IP Reader
An asynchronous text processor and batch file reader.

## Table of contents
  - [How to use the program](#How-to-use-the-program)
    - [Run](#RUN)
    - [Test](#test)
	- [Profile](#profile)
  - [Results from processing the 100G file](#Results-from-processing-the-100G-file)
  - [Architecture description](#Architecture-description)
    - [ipreader package](#ipreader-package)
	- [ipcounter package](#ipcounter-package)
  - [Profiling](#Profiling)
    - [Memory profiling results](#Memory-profiling-results)

## How to use the program
Simply open the `main.go` file, and change the value of the variable `file, err := os.Open("/run/media/jenia/My Book/ip_addresses")` with the path to the IP file you want to process

### RUN

```
go run .
```

### Test
```
go test ./...
```
### Profile
```
go test -bench=. -benchmem -memprofile mem.prof -cpuprofile cpu.prof -benchtime=60s
go tool pprof mem.prof
top 20
```

## Results from processing the 100G file

```
$ time go run .
Count is: 301774584

real	1m10.062s
user	1m11.607s
sys	0m4.389s
```

Sanity check:
```
$ sort "/run/media/jenia/My Book/ip_addresses" | uniq | wc -l
301774584
```

## Architecture description

This program reads a file that contains one IP address per line and counts the number of unique IPs contained in that file. Reading and counting happens asynchronously.

This program has two packages:
- ipreader
- ipcounter

This program is carefully designed to process large volume of data while using a minimum of resources:

- This program is profiled using Golang pprof tooling
- ips are treated as numbers rather than strings to optimize memory and lookup performance
- ips are stored in a kind of multiplication table to optimize memory usage
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

Using an interface, I decouple `ReadFile` function from whatever implements the `Counter` interface which is Golang best practices. In this case it's `ipcounter.IpCounter` that implements the `Counter` interface by asynchronously counting unique IPs in the slice `AddIpSlice(ips []uint32)`.
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

This program is meticulously profiled using Golang pprof tooling

### Memory profiling results
```
>go test -bench=. -benchmem -memprofile mem.prof -cpuprofile cpu.prof -benchtime=60s
go tool pprof mem.prof
top 20
   11264MB 50.36% 50.36%    11264MB 50.36%  github.com/Ecwid/new-job/ipcounter.NewIpCounter (inline)
 3558.05MB 15.91% 66.27%  3558.05MB 15.91%  net.ParseIP (inline)
 3542.05MB 15.84% 82.11%  3542.05MB 15.84%  bufio.(*Scanner).Text (inline)
 2255.24MB 10.08% 92.19%  9355.35MB 41.83%  github.com/Ecwid/new-job/ipreader.ReadFile
  627.40MB  2.81% 95.00%   627.40MB  2.81%  bytes.growSlice
  451.51MB  2.02% 97.02%  1739.42MB  7.78%  github.com/Ecwid/new-job.writeIpsToFile
  335.51MB  1.50% 98.52%   659.51MB  2.95%  fmt.Sprintf
     322MB  1.44%   100%      322MB  1.44%  net/netip.Addr.string4 (inline)
       1MB 0.0045%   100%   628.40MB  2.81%  bytes.(*Buffer).grow
```
