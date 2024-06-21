package ipcounter

import (
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
