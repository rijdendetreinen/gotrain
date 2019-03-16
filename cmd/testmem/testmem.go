package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/rijdendetreinen/gotrain/models"

	"github.com/rijdendetreinen/gotrain/parsers"
)

func main() {
	filename := os.Args[1]
	total := 10000

	// PrintMemUsage(0)
	start := time.Now()

	var services []models.Service

	for i := 0; i < total; i++ {
		// fmt.Printf("Opening %s\n", filename)

		f, err := os.Open(filename)

		if err != nil {
			panic(err)
		}

		services = append(services, parsers.ParseRitMessage(f))

		f.Close()

		if i%250 == 0 {
			PrintMemUsage(i)
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)

	msgPerSec := (float64(total) / elapsed.Seconds())

	PrintMemUsage(total)

	fmt.Printf("Time: %s, Messages per second: %.2f", elapsed, msgPerSec)
	fmt.Println()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage(counter int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Count = %v", counter)
	fmt.Printf("\tAlloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
