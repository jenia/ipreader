package ipcounter

import (
	"sync"
)

const (
	ipSize = 32
)

type IpCounter struct {
	ipSlices chan []uint32
	Counter  uint32
	ips      []uint32
	mtx      sync.Mutex
	closed   bool
}

// Can only create one of these objects
func NewIpCounter() *IpCounter {
	ipSlices := make(chan []uint32, 1)
	// I'm making a multiplication table of sorts to hold the IPs
	// the number 134_217_728 = number of IP v4 / 32
	// 32 is the length of uint32
	// So, 134_217_728 * 32 is the dimension of the table
	ips := make([]uint32, 134_217_728)
	return &IpCounter{ipSlices: ipSlices, ips: ips}
}

func (i *IpCounter) Count(wg *sync.WaitGroup) {
	defer wg.Done()
	for ipSlice := range i.ipSlices {
		for _, ip := range ipSlice {
			i.mtx.Lock()
			q := ip / ipSize
			r := ip % ipSize
			row := i.ips[q]
			if row &(1<<r) == 0 {
				row |= 1 << r
				i.ips[q] = row
				i.Counter++
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
