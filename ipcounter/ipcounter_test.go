package ipcounter

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"
)

// I'm not using table driven tests because the logic with the go-routines, wait groups and so on, makes
// table driven tests a little less readable in this particular case.
func TestProcessIps(t *testing.T) {
	t.Run("Given 3 different IPs, when channel closed, then function returns and counter is 3", func(t *testing.T) {
		// setup
		ips := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
		wg := sync.WaitGroup{}
		ipCounter := NewIpCounter()

		// test
		wg.Add(1)
		go ipCounter.Count(&wg)
		ipCounter.AddIpSlice(ips)
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

		// test
		wg.Add(1)
		go ipCounter.Count(&wg)
		ipCounter.AddIpSlice(ips)
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
	b.ResetTimer()
	buf := make([]byte, 4)
	ipSlice := make([]string, 100)
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint32(buf, rand.Uint32())
		ip := fmt.Sprintf("%s", net.IP(buf))
		if i%100 == 0 {
			c.AddIpSlice(ipSlice)
		} else {
			ipSlice[i%100] = ip
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
	buf := make([]byte, 4)
	ipSlice := make([]string, 1000)
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint32(buf, rand.Uint32())
		ip := fmt.Sprintf("%s", net.IP(buf))
		if i%1000 == 0 {
			c.AddIpSlice(ipSlice)
		} else {
			ipSlice[i%1000] = ip
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
	b.ResetTimer()
	buf := make([]byte, 4)
	ipSlice := make([]string, 100)
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint32(buf, rand.Uint32())
		ip := fmt.Sprintf("%s", net.IP(buf))
		if i%100 == 0 {
			c.AddIpSlice(ipSlice)
		} else {
			ipSlice[i%100] = ip
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
	b.ResetTimer()
	buf := make([]byte, 4)
	ipSlice := make([]string, 1000)
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint32(buf, rand.Uint32())
		ip := fmt.Sprintf("%s", net.IP(buf))
		if i%1000 == 0 {
			c.AddIpSlice(ipSlice)
		} else {
			ipSlice[i%1000] = ip
		}
	}
	c.Close()
	wg.Wait()
}
