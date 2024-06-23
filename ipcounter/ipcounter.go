package ipcounter

import (
	"sync"
)

type IpCounter struct {
	ipSlices chan []string
	Counter  uint64
	ipMap    map[string]bool
	mtx      sync.Mutex
	closed   bool
}

// Can only create one of these objects
func NewIpCounter() *IpCounter {
	ipSlices := make(chan []string, 1)
	ipMap := make(map[string]bool, 1)
	return &IpCounter{ipSlices: ipSlices, ipMap: ipMap}
}

func (i *IpCounter) Count(wg *sync.WaitGroup) {
	defer wg.Done()
	for ipSlice := range i.ipSlices {
		for _, ip := range ipSlice {
			i.mtx.Lock()
			if _, ok := i.ipMap[ip]; !ok {
				i.Counter++
				i.ipMap[ip] = true
			}
			i.mtx.Unlock()
		}
	}
}

// Not thread safe
func (i *IpCounter) Close() {
	if i.closed == true {
		return
	}
	close(i.ipSlices)
	i.closed = true
}

func (i *IpCounter) AddIpSlice(ips []string) {
	i.ipSlices <- ips
}
