package main

import (
	"fmt"
	"github.com/Ecwid/new-job/ipcounter"
	"github.com/Ecwid/new-job/ipreader"
	"os"
	"sync"
)

func main() {
	file, err := os.Open("./ips.txt")
	if err != nil {
		fmt.Printf("open file: %s", err.Error())
		panic(err)
	}
	defer file.Close()

	wg := &sync.WaitGroup{}
	ipCounter := ipcounter.NewIpCounter()
	wg.Add(1)
	go ipCounter.Count(wg)
	buf := make([]byte, 2^20)
	ipreader.ReadFile(file, ipCounter, buf)
	wg.Wait()
	fmt.Printf("Count is: %d\n", ipCounter.Counter)
}
