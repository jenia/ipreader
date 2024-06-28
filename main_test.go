package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Ecwid/new-job/ipcounter"
	"github.com/Ecwid/new-job/ipreader"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
)

const fileName = "ipsfortesting.txt"

func BenchmarkTestEntireProgram(b *testing.B) {
	writeIpsToFile()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("open file: %s", err.Error())
			panic(err)
		}
		defer file.Close()

		wg := &sync.WaitGroup{}
		ipCounter := ipcounter.NewIpCounter()
		wg.Add(1)
		go ipCounter.Count(wg)
		buf := make([]byte, readBufferSize)
		ipreader.ReadFile(file, ipCounter, buf)
		wg.Wait()
		fmt.Printf("Count is: %d\n", ipCounter.Counter)
	}
}

func writeIpsToFile() {
	const (
		giga = 1024 * 1024 * 10
	)

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	ipBuf := make([]byte, 4)
	fileBuf := &bytes.Buffer{}
	for i := range giga {
		binary.LittleEndian.PutUint32(ipBuf, rand.Uint32())
		_, err = fileBuf.WriteString(fmt.Sprintf("%s\n", net.IP(ipBuf)))
		if err != nil {
			panic(err)
		}
		if i%1024 == 0 {
			file.Write(fileBuf.Bytes())
			fileBuf = &bytes.Buffer{}
		}
	}
}
