package ipreader

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

type Counter interface {
	AddIpSlice(ips []string)
	Close()
}

// TODO: write error if IP is not parsable
func ReadFile(file io.Reader, counter Counter, buffer []byte) error {
	if len(buffer) < 17 {
		return errors.New("buffer size below length of ipv4 size")
	}

	defer counter.Close()
	var tmpBuffer bytes.Buffer

	var residual []byte
	for {
		bytesRead, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			fmt.Printf("read file: %s\n", err.Error())
			return fmt.Errorf("read file: %w", err)
		}
		if bytesRead == 0 {
			break
		}

		tmpBuffer.Write(buffer[:bytesRead])
		scanner := bufio.NewScanner(&tmpBuffer)
		scanner.Split(bufio.ScanLines)

		var ips []string
		for scanner.Scan() {
			ip := scanner.Text()
			if ip == "" {
				continue
			}
			ips = append(ips, ip)
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("scan buffer for new lines: %s\n", err.Error())
			return fmt.Errorf("scanner buffer for new lines: %w", err)
		}

		tmpBuffer.Reset()
		if bytesRead == len(buffer) {
			residual = keepResidual(ips)
			if len(residual) > 0 {
				ips = ips[:len(ips)-1]
				tmpBuffer.Write(residual)
			}
		}
		counter.AddIpSlice(ips)
	}
	return nil
}

func keepResidual(lines []string) []byte {
	if len(lines) > 0 {
		lastLine := lines[len(lines)-1]
		if len(lastLine) > 0 && lastLine[len(lastLine)-1] != '\n' {
			return []byte(lastLine)
		}
	}
	return []byte{}
}
