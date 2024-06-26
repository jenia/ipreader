package ipreader

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type Counter interface {
	AddIpSlice(ips []uint32)
	Close()
}

func Ipv4ToInt(ipaddr net.IP) (uint32, error) {
	if ipaddr.To4() == nil {
		return 0, errors.New(fmt.Sprintf("not an ip4 addr: %+v", ipaddr))
	}
	return binary.BigEndian.Uint32(ipaddr.To4()), nil
}

func ReadFile(file io.Reader, counter Counter, buffer []byte) error {
	defer counter.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(buffer, 32)
	scanner.Split(bufio.ScanLines)

	var ips []uint32
	for scanner.Scan() {
		ip := scanner.Text()
		if ip == "" {
			fmt.Print("error: empty line")
			continue
		}
		ipInt, err := Ipv4ToInt(net.ParseIP(ip))
		if err != nil {
			fmt.Printf("convert ip to uint32: %s", err.Error())
		}
		ips = append(ips, ipInt)
		if len(ips) >= 100 {
			counter.AddIpSlice(ips)
			ips = nil
		}
	}
	if len(ips) > 0 {
		counter.AddIpSlice(ips)
	}

	// TODO: properly check for errors in the scanner. I think we might need to check it at each iteration?
	if err := scanner.Err(); err != nil {
		fmt.Printf("scan buffer for new lines: %s\n", err.Error())
		return fmt.Errorf("scan buffer for new lines: %w", err)
	}

	return nil
}
