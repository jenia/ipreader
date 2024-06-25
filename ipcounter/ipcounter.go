package ipcounter

import (
	"sync"
)

type IpCounter struct {
	ipSlices chan []uint32
	Counter  uint64
	ips      []bool
	mtx      sync.Mutex
	closed   bool
}

// Can only create one of these objects
func NewIpCounter() *IpCounter {
	ipSlices := make(chan []uint32, 1)
	ips := make([]bool, uint64(1) << 32)
	return &IpCounter{ipSlices: ipSlices, ips: ips}
}

func (i *IpCounter) Count(wg *sync.WaitGroup) {
	defer wg.Done()
	for ipSlice := range i.ipSlices {
		for _, ip := range ipSlice {
			i.mtx.Lock()
			if i.ips[ip] == false {
				i.Counter++
				i.ips[ip] = true
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

func (i *IpCounter) AddIpSlice(ips []uint32) {
	i.ipSlices <- ips
}
