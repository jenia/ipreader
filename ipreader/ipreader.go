package ipreader

import (
	"bufio"
	"fmt"
	"io"
)

type Counter interface {
	AddIpSlice(ips []string)
	Close()
}

// TODO: write error if IP is not parsable
func ReadFile(file io.Reader, counter Counter, buffer []byte) error {
	defer counter.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(buffer, 32)
	scanner.Split(bufio.ScanLines)

	var ips []string
	for scanner.Scan() {
		ip := scanner.Text()
		if ip == "" {
			// TODO: this is an error actually... not fatal but should be reported
			continue
		}
		ips = append(ips, ip)
		if len(ips) >= 100 {
			counter.AddIpSlice(ips)
			ips = nil // Clear the slice to free memory
		}
	}
	if len(ips) > 0 {
		counter.AddIpSlice(ips)
	}

	// Check for errors in the scanner
	if err := scanner.Err(); err != nil {
		fmt.Printf("scan buffer for new lines: %s\n", err.Error())
		return fmt.Errorf("scan buffer for new lines: %w", err)
	}

	return nil
} 
