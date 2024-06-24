# IpReader
An asynchronous text processor and batch file reader.

## Limitations
This program consumes a lot of memory because it holds the IPs in memory in a hash table.
Normally, this program should search/write the IPs to a relational database:
```
BEGIN SERIALIZABLE;
SELECT ip
FROM ips
WHERE NOT EXISTS (
    SELECT 1
    FROM (VALUES
        (ip1),
        (ip2),
    ) AS batch(ip)
    WHERE ips.ipv4 = batch.ip
);
INSERT INTO ips (ipv4) values ip1, ip2,...
ON CONFLICT DO NOTHING;
COMMIT;
```

Ofcourse the table `ips` needs to be indexed. If we're in a multi-tenant environment, then we need to search/insert using the column `tenant_id` (or something similar) as well.

I decided not to use a database for demo purposes but instead focus on Golang and its asynchronous functionality.

## Architecture description

This program reads a file that contains one IP address per line and counts the number of unique IP.

This program has two packages:
- ipreader
- ipcounter

### ipreader package
This package defines a function `ReadFile` that reads a file into a buffer and sends all the IPs contained in that buffer to `Counter.AddIpSlice(ips []string)`.

Note the signature of the function `func ReadFile(file io.Reader, counter Counter, buffer []byte) error`

It takes
1. `io.Reader` file handle
2. `Counter` interface
```
type Counter interface {
	AddIpSlice(ips []string)
	Close()
}
```
3. a `[]byte` buffer

Using an interface, I decouple `ReadFile` function from whatever implements the `Counter` interface which is Golang best practices. In this case it's `ipcounter.IpCounter` that implements the `Counter` interface by asynchronously counting the IPs in the slice `AddIpSlice(ips []string)`.
This decoupling allows for higher modularity which has important implications:
1. Simpler code
2. Simpler testing. Each feature can be unit tested separately. Same thing for benchmarking. See the tests and benchmarks for details.

Furthermore, using a buffer + goroutines we can bulk read and async process the IPs.

### ipcounter package
This package defines a struct `IpCounter` that counts unique IPs sent to it using the function `AddIpSlice`. Each item in the slice is one IP.
The struct `IpCounter` is meant to be instantiated once.
The function `Count` can run, and it meant to be run (see benchmarking), on multiple go-routines

See the main function or the bench marking tests in `ipcounter` package to see how to instantiate and use `IpCounter`

## CPU Benchmarking results.

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

