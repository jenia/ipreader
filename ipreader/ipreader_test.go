package ipreader

import (
	"strings"
	"testing"
)

type counterMock struct {
	ipsCh chan []string
}

func (c *counterMock) AddIpSlice(ips []string) {
	c.ipsCh <- ips
}

func (c *counterMock) Close() {
	close(c.ipsCh)
}

func TestRead(t *testing.T) {
	t.Run("Given file has 5 IPs, when buffer can hold part of file, then return IPs in multiple channels", func(t *testing.T) {
		//setup
		fileContent := `145.67.23.4
8.34.5.23
89.54.3.124
89.54.3.124
3.45.71.5
`
		r := strings.NewReader(fileContent)
		ipChan := make(chan []string, 1)
		sizeOfIP4Address := 20
		buffer := make([]byte, sizeOfIP4Address)
		counterMock := counterMock{ipsCh: ipChan}

		// test
		go ReadFile(r, &counterMock, buffer)

		// verify
		i := 0
		lines := strings.Split(fileContent, "\n")
		lines = lines[:len(lines)-1] // drop the line with only the \n
		for ipSlice := range ipChan {
			for j, ip := range ipSlice {
				if lines[i] != ipSlice[j] {
					t.Fatalf("ip should be %s, but is %s", lines[i], ip)
				}
				i++
			}
		}
		if i != len(lines) {
			t.Errorf("file has 5 ips, but channel received only %d ips", i)
		}
	})

	t.Run("Given file has 5 IPs, when buffer can hold entire file, then channel has 1 slice with 5 IPs", func(t *testing.T) {
		//setup
		fileContent := `145.67.23.4
8.34.5.23
89.54.3.124
89.54.3.124
3.45.71.5
`
		r := strings.NewReader(fileContent)
		ipChan := make(chan []string, 1)
		buffer := make([]byte, 1024)
		counterMock := counterMock{ipsCh: ipChan}

		// test
		go ReadFile(r, &counterMock, buffer)

		// verify
		lines := strings.Split(fileContent, "\n")
		lines = lines[:len(lines)-1] // drop the line with only the \n
		ipSlice := <-ipChan
		if len(ipSlice) != len(lines) {
			t.Fatalf("ipSlice len should %d, but instead is %d: %+v", len(lines), len(ipSlice), ipSlice)
		}
		for i, line := range lines {
			if line != ipSlice[i] {
				t.Fatalf("ip should be %s, but is %s", line, ipSlice[i])
			}
		}
	})

	t.Run("Given file has 5 IPs, when file last character is not a newline, then return 5 ips", func(t *testing.T) {
		//setup
		fileContent := `145.67.23.4
8.34.5.23
89.54.3.124
89.54.3.124
3.45.71.5`
		r := strings.NewReader(fileContent)
		ipChan := make(chan []string, 1)
		buffer := make([]byte, 1024)
		counterMock := counterMock{ipsCh: ipChan}

		// test
		go ReadFile(r, &counterMock, buffer)

		// verify
		lines := strings.Split(fileContent, "\n")
		ipSlice := <-ipChan
		if len(ipSlice) != len(lines) {
			t.Fatalf("ipSlice len should %d, but instead is %d: %+v", len(lines), len(ipSlice), ipSlice)
		}
		for i, line := range lines {
			if line != ipSlice[i] {
				t.Fatalf("ip should be %s, but is %s", line, ipSlice[i])
			}
		}
	})

	t.Run("Given file is empty, then return 0 ips", func(t *testing.T) {
		//setup
		fileContent := ``
		r := strings.NewReader(fileContent)
		ipChan := make(chan []string, 1)
		buffer := make([]byte, 1024)
		counterMock := counterMock{ipsCh: ipChan}

		// test
		go ReadFile(r, &counterMock, buffer)

		// verify
		ipSlice := <-ipChan
		if len(ipSlice) != 0 {
			t.Fatalf("ipSlice len should 0, but instead is %d: %+v", len(ipSlice), ipSlice)
		}
	})

	t.Run("Given file has only new line, then return 0 ips", func(t *testing.T) {
		//setup
		fileContent := `
`
		r := strings.NewReader(fileContent)
		ipChan := make(chan []string, 1)
		buffer := make([]byte, 1024)
		counterMock := counterMock{ipsCh: ipChan}

		// test
		go ReadFile(r, &counterMock, buffer)

		// verify
		ipSlice := <-ipChan
		if len(ipSlice) != 0 {
			t.Fatalf("ipSlice len should 0, but instead is %d: %+v", len(ipSlice), ipSlice)
		}
	})
}
