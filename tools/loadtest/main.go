// Command loadtest mengirim beban HTTP sederhana ke satu endpoint.
//
// Contoh:
//
//	go run ./tools/loadtest -url http://localhost:3000/api/v1/health -n 2000 -c 50
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	url := flag.String("url", "http://localhost:3000/api/v1/health", "target URL")
	total := flag.Int("n", 1000, "jumlah request")
	concurrency := flag.Int("c", 25, "jumlah worker paralel")
	timeout := flag.Duration("timeout", 10*time.Second, "timeout per request")
	flag.Parse()

	if *total < 1 || *concurrency < 1 {
		fmt.Fprintln(os.Stderr, "n dan c harus lebih besar dari 0")
		os.Exit(1)
	}

	client := &http.Client{Timeout: *timeout}
	var ok, failed int64
	latencies := make([]time.Duration, 0, *total)
	var mu sync.Mutex

	jobs := make(chan struct{}, *total)
	for i := 0; i < *total; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				began := time.Now()
				resp, err := client.Get(*url)
				elapsed := time.Since(began)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					continue
				}
				_ = resp.Body.Close()

				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					atomic.AddInt64(&ok, 1)
				} else {
					atomic.AddInt64(&failed, 1)
				}

				mu.Lock()
				latencies = append(latencies, elapsed)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	fmt.Printf("target      : %s\n", *url)
	fmt.Printf("request     : %d (paralel %d)\n", *total, *concurrency)
	fmt.Printf("sukses/gagal: %d / %d\n", ok, failed)
	fmt.Printf("durasi      : %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("throughput  : %.1f req/s\n", float64(*total)/elapsed.Seconds())
	if len(latencies) > 0 {
		fmt.Printf("latensi p50 : %s\n", percentile(latencies, 50).Round(time.Microsecond))
		fmt.Printf("latensi p95 : %s\n", percentile(latencies, 95).Round(time.Microsecond))
		fmt.Printf("latensi p99 : %s\n", percentile(latencies, 99).Round(time.Microsecond))
	}
}

func percentile(sorted []time.Duration, p int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	index := (len(sorted) - 1) * p / 100
	return sorted[index]
}
