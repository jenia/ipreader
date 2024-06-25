package ipcounter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"
)

func Ipv4ToInt(ipaddr net.IP) (uint32, error) {
	if ipaddr.To4() == nil {
		return 0, errors.New(fmt.Sprintf("not an ip addr: %+v", ipaddr))
	}
	return binary.BigEndian.Uint32(ipaddr.To4()), nil
}

func createIPInts(ips []string) ([]uint32, error) {
	ipInts := make([]uint32, len(ips))
	for i, ip0 := range ips {
		ipInt, err := Ipv4ToInt(net.ParseIP(ip0))
		if err != nil {
			return []uint32{}, fmt.Errorf("convert ip to uint32: %w", err)
		}
		ipInts[i] = ipInt
	}
	return ipInts, nil
}

// I'm not using table driven tests because the logic with the go-routines, wait groups and so on, makes
// table driven tests a little less readable in this particular case.
func TestProcessIps(t *testing.T) {
	t.Run("Given 3 different IPs, when channel closed, then function returns and counter is 3", func(t *testing.T) {
		// setup
		ips := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
		wg := sync.WaitGroup{}
		ipCounter := NewIpCounter()
		ipInts, err := createIPInts(ips)
		if err != nil {
			t.Fatalf("create IPS: %s", err.Error())
		}

		// test
		wg.Add(1)
		go ipCounter.Count(&wg)
		ipCounter.AddIpSlice(ipInts)
		ipCounter.Close()
		wg.Wait()

		// verify
		if ipCounter.Counter != 3 {
			t.Errorf("Counter should be 3 but is %d", ipCounter.Counter)
		}
	})

	t.Run("Given 3 different IPs out of 4, when channel closed, then function returns and counter is 3", func(t *testing.T) {
		// setup
		ips := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.3"}
		wg := sync.WaitGroup{}
		ipCounter := NewIpCounter()
		ipInts, err := createIPInts(ips)
		if err != nil {
			t.Fatalf("create IPS: %s", err.Error())
		}


		// test
		wg.Add(1)
		go ipCounter.Count(&wg)
		ipCounter.AddIpSlice(ipInts)
		ipCounter.Close()
		wg.Wait()

		// verify
		if ipCounter.Counter != 3 {
			t.Errorf("Counter should be 3 but is %d", ipCounter.Counter)
		}
	})
}

func BenchmarkCount1GoRoutine100ItemSlice(b *testing.B) {
	c := NewIpCounter()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go c.Count(&wg)
	ipSlice := make([]uint32, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%100 == 0 {
			c.AddIpSlice(ipSlice)
		} else {
			ipSlice[i%100] = rand.Uint32()
		}
	}
	c.Close()
	wg.Wait()
}

func BenchmarkCount1GoRoutine1000ItemSlice(b *testing.B) {
	c := NewIpCounter()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go c.Count(&wg)
	b.ResetTimer()
	ipSlice := make([]uint32, 1000)
	for i := 0; i < b.N; i++ {
		if i%1000 == 0 {
			c.AddIpSlice(ipSlice)
		} else {
			ipSlice[i%1000] = rand.Uint32()
		}
	}
	c.Close()
	wg.Wait()
}

func BenchmarkCount10GoRoutine100ItemSlice(b *testing.B) {
	c := NewIpCounter()
	wg := sync.WaitGroup{}
	wg.Add(10)
	for range 10 {
		go c.Count(&wg)
	}
	ipSlice := make([]uint32, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%100 == 0 {
			c.AddIpSlice(ipSlice)
		} else {
			ipSlice[i%100] = rand.Uint32()
		}
	}
	c.Close()
	wg.Wait()
}

func BenchmarkCount10GoRoutine1000ItemSlice(b *testing.B) {
	c := NewIpCounter()
	wg := sync.WaitGroup{}
	wg.Add(10)
	for range 10 {
		go c.Count(&wg)
	}
	ipSlice := make([]uint32, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%1000 == 0 {
			c.AddIpSlice(ipSlice)
		} else {
			ipSlice[i%1000] = rand.Uint32()
		}
	}
	c.Close()
	wg.Wait()
}
