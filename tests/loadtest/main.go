package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBlogURL = "http://localhost:8090"
	userUUID       = "11111111-1111-1111-1111-111111111111"
	requestTimeout = 5 * time.Second

	createGoroutines = 100
	listGoroutines   = 1000
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type phaseStats struct {
	success     atomic.Int64
	failed      atomic.Int64
	totalTimeNs atomic.Int64
}

func (s *phaseStats) record(duration time.Duration, err error) {
	if err != nil {
		s.failed.Add(1)
	} else {
		s.success.Add(1)
		s.totalTimeNs.Add(duration.Nanoseconds())
	}
}

func (s *phaseStats) print(phaseName string) {
	succ := s.success.Load()
	fail := s.failed.Load()
	total := succ + fail
	var pct float64
	if total > 0 {
		pct = float64(succ) / float64(total) * 100
	}
	var avgMs float64
	if succ > 0 {
		avgMs = float64(s.totalTimeNs.Load()) / float64(succ) / 1e6
	}

	fmt.Printf("\n=== %s ===\n", phaseName)
	fmt.Printf("Total requests:       %d\n", total)
	fmt.Printf("Successful:           %d\n", succ)
	fmt.Printf("Failed (timeout/err): %d\n", fail)
	fmt.Printf("Success rate:         %.1f%%\n", pct)
	fmt.Printf("Avg response time:    %.2f ms\n", avgMs)
}

func main() {
	blogURL := getEnv("BLOG_URL", defaultBlogURL)
	client := &http.Client{Timeout: requestTimeout}

	fmt.Printf("Load test client started\n")
	fmt.Printf("Target: %s\n", blogURL)
	fmt.Printf("User UUID: %s\n", userUUID)

	////////////////////////////////////////////////////////////////////////////

	fmt.Printf("\n--- Phase 1: Creating posts (%d goroutines) ---\n", createGoroutines)
	createStats := &phaseStats{}
	var wg sync.WaitGroup

	for i := 0; i < createGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			body, _ := json.Marshal(map[string]string{
				"contentText": fmt.Sprintf("Load test post #%d at %s", idx, time.Now().Format(time.RFC3339Nano)),
			})

			req, _ := http.NewRequest("POST", blogURL+"/v1/posts", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-Uuid", userUUID)

			start := time.Now()
			resp, err := client.Do(req)
			dur := time.Since(start)

			if err != nil {
				createStats.record(dur, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				createStats.record(dur, fmt.Errorf("status %d", resp.StatusCode))
				return
			}

			createStats.record(dur, nil)
		}(i)
	}

	wg.Wait()
	createStats.print("Phase 1 completed")

	////////////////////////////////////////////////////////////////////////////

	fmt.Printf("\n--- Phase 2: Listing posts (%d goroutines) ---\n", listGoroutines)
	listStats := &phaseStats{}

	for i := 0; i < listGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			limit := rand.Intn(20) + 1
			url := fmt.Sprintf("%s/v1/posts?limit=%d", blogURL, limit)

			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("X-User-Uuid", userUUID)

			start := time.Now()
			resp, err := client.Do(req)
			dur := time.Since(start)

			if err != nil {
				listStats.record(dur, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				listStats.record(dur, fmt.Errorf("status %d", resp.StatusCode))
				return
			}

			listStats.record(dur, nil)
		}()
	}

	wg.Wait()
	listStats.print("Phase 2 completed")
}
