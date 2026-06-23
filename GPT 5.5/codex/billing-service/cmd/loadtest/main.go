// Command loadtest runs a bounded, read-only HTTP load test against Billing Service.
// It authenticates once, then concurrently invokes a protected endpoint and prints
// latency percentiles and HTTP status counts. No billing data is created or changed.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

type sample struct {
	Latency time.Duration
	Status  int
	Err     error
}

func main() {
	baseURL := flag.String("base-url", env("LOAD_BASE_URL", "http://localhost:8080"), "service base URL")
	endpoint := flag.String("endpoint", env("LOAD_ENDPOINT", "/api/v1/plans"), "protected GET endpoint")
	requests := flag.Int("requests", envInt("LOAD_REQUESTS", 1000), "total requests")
	concurrency := flag.Int("concurrency", envInt("LOAD_CONCURRENCY", 25), "parallel workers")
	email := flag.String("email", env("LOAD_EMAIL", "admin@example.com"), "login email")
	password := flag.String("password", env("LOAD_PASSWORD", "admin-password-change-me"), "login password")
	timeout := flag.Duration("timeout", envDuration("LOAD_TIMEOUT", 10*time.Second), "per-request timeout")
	flag.Parse()

	if *requests < 1 || *concurrency < 1 || *concurrency > *requests || *timeout <= 0 {
		fmt.Fprintln(os.Stderr, "requests and concurrency must be positive, concurrency must not exceed requests, and timeout must be positive")
		os.Exit(2)
	}

	base := strings.TrimRight(*baseURL, "/")
	client := &http.Client{Transport: &http.Transport{MaxIdleConns: *concurrency * 2, MaxIdleConnsPerHost: *concurrency * 2, IdleConnTimeout: 30 * time.Second}}
	token, err := login(client, base, *email, *password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		os.Exit(1)
	}

	jobs := make(chan struct{})
	results := make(chan sample, *requests)
	var wg sync.WaitGroup
	for range *concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				results <- request(client, base+*endpoint, token, *timeout)
			}
		}()
	}

	started := time.Now()
	for range *requests {
		jobs <- struct{}{}
	}
	close(jobs)
	wg.Wait()
	close(results)
	elapsed := time.Since(started)

	latencies := make([]time.Duration, 0, *requests)
	statuses := map[int]int{}
	var failures atomic.Int64
	for result := range results {
		latencies = append(latencies, result.Latency)
		statuses[result.Status]++
		if result.Err != nil || result.Status < 200 || result.Status >= 300 {
			failures.Add(1)
		}
	}
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	fmt.Printf("Load test: GET %s\n", base+*endpoint)
	fmt.Printf("Requests: %d | Concurrency: %d | Duration: %s | Throughput: %.2f req/s\n", *requests, *concurrency, elapsed.Round(time.Millisecond), float64(*requests)/elapsed.Seconds())
	fmt.Printf("Success: %d | Failed: %d | Statuses: %s\n", *requests-int(failures.Load()), failures.Load(), formatStatuses(statuses))
	fmt.Printf("Latency: min=%s avg=%s p50=%s p95=%s p99=%s max=%s\n", latencies[0].Round(time.Microsecond), average(latencies).Round(time.Microsecond), percentile(latencies, .50).Round(time.Microsecond), percentile(latencies, .95).Round(time.Microsecond), percentile(latencies, .99).Round(time.Microsecond), latencies[len(latencies)-1].Round(time.Microsecond))

	if failures.Load() > 0 {
		os.Exit(1)
	}
}

func login(client *http.Client, base, email, password string) (string, error) {
	body, _ := json.Marshal(map[string]string{"email": email, "password": password})
	req, err := http.NewRequest(http.MethodPost, base+"/api/v1/auth/login", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return "", fmt.Errorf("HTTP %d: %s", response.StatusCode, strings.TrimSpace(string(b)))
	}
	var payload loginResponse
	if err = json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.AccessToken == "" {
		return "", fmt.Errorf("response did not contain an access token")
	}
	return payload.AccessToken, nil
}

func request(client *http.Client, url, token string, timeout time.Duration) sample {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return sample{Err: err}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	started := time.Now()
	response, err := client.Do(req)
	latency := time.Since(started)
	if err != nil {
		return sample{Latency: latency, Err: err}
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, response.Body)
	return sample{Latency: latency, Status: response.StatusCode}
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	return sorted[int(math.Ceil(p*float64(len(sorted))))-1]
}
func average(values []time.Duration) time.Duration {
	var sum time.Duration
	for _, value := range values {
		sum += value
	}
	return sum / time.Duration(len(values))
}
func formatStatuses(statuses map[int]int) string {
	keys := make([]int, 0, len(statuses))
	for key := range statuses {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%d=%d", key, statuses[key]))
	}
	return strings.Join(parts, ", ")
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err == nil && value > 0 {
		return value
	}
	return fallback
}
func envDuration(key string, fallback time.Duration) time.Duration {
	value, err := time.ParseDuration(os.Getenv(key))
	if err == nil && value > 0 {
		return value
	}
	return fallback
}
