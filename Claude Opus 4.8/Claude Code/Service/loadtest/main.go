// Command loadtest is a small self-contained HTTP load generator for the
// billing service. It logs in, seeds the data it needs, then drives a pool of
// concurrent workers against a chosen scenario and reports latency percentiles
// and throughput. Standard library only.
//
// Usage:
//
//	go run ./loadtest -url http://localhost:8080 -scenario mixed -c 50 -d 20s
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type config struct {
	baseURL    string
	scenario   string
	concurrency int
	duration   time.Duration
	adminEmail string
	adminPass  string
}

type result struct {
	latency time.Duration
	status  int
	err     bool
}

type seed struct {
	token      string
	customerID string
	planID     string
	subID      string
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.baseURL, "url", "http://localhost:8080", "base URL")
	flag.StringVar(&cfg.scenario, "scenario", "mixed", "read | usage | invoice | mixed")
	flag.IntVar(&cfg.concurrency, "c", 50, "concurrent workers")
	flag.DurationVar(&cfg.duration, "d", 15*time.Second, "test duration")
	flag.StringVar(&cfg.adminEmail, "email", "admin@billing.local", "admin email")
	flag.StringVar(&cfg.adminPass, "pass", "admin12345", "admin password")
	flag.Parse()

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        2000,
			MaxIdleConnsPerHost: 2000,
			MaxConnsPerHost:     0,
			IdleConnTimeout:     60 * time.Second,
		},
	}

	s, err := setup(client, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("scenario=%s concurrency=%d duration=%s\n", cfg.scenario, cfg.concurrency, cfg.duration)
	run(client, cfg, s)
}

func setup(client *http.Client, cfg config) (seed, error) {
	var s seed
	// Login.
	body := map[string]string{"email": cfg.adminEmail, "password": cfg.adminPass}
	resp, code, err := doJSON(client, "POST", cfg.baseURL+"/api/v1/auth/login", "", body)
	if err != nil || code != 200 {
		return s, fmt.Errorf("login: code=%d err=%v", code, err)
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.Unmarshal(resp, &tok)
	s.token = tok.AccessToken
	if s.token == "" {
		return s, fmt.Errorf("no access token")
	}

	// Customer.
	resp, code, err = doJSON(client, "POST", cfg.baseURL+"/api/v1/customers", s.token,
		map[string]any{"email": "load@example.com", "name": "Load", "country": "US", "currency": "USD"})
	if err != nil || code >= 300 {
		return s, fmt.Errorf("create customer: code=%d err=%v", code, err)
	}
	s.customerID = idOf(resp)

	// Plan (usage-based so invoice generation has work to do).
	resp, code, err = doJSON(client, "POST", cfg.baseURL+"/api/v1/plans", s.token,
		map[string]any{"name": "Load Plan", "currency": "USD", "billing_interval": "monthly",
			"pricing_model": "usage_based", "base_price_cents": 1000, "included_units": 0,
			"unit_price_cents": 2, "usage_metric": "api_calls"})
	if err != nil || code >= 300 {
		return s, fmt.Errorf("create plan: code=%d err=%v", code, err)
	}
	s.planID = idOf(resp)

	// Subscription.
	resp, code, err = doJSON(client, "POST", cfg.baseURL+"/api/v1/subscriptions", s.token,
		map[string]any{"customer_id": s.customerID, "plan_id": s.planID, "quantity": 1})
	if err != nil || code >= 300 {
		return s, fmt.Errorf("create subscription: code=%d err=%v", code, err)
	}
	s.subID = idOf(resp)

	// Tax rule so invoices exercise tax.
	_, _, _ = doJSON(client, "POST", cfg.baseURL+"/api/v1/tax-rules", s.token,
		map[string]any{"country": "US", "region": "", "tax_name": "Tax", "tax_rate_basis_points": 2000})

	return s, nil
}

func run(client *http.Client, cfg config, s seed) {
	deadline := time.Now().Add(cfg.duration)
	var wg sync.WaitGroup
	perWorker := make([][]result, cfg.concurrency)
	var counter int64

	start := time.Now()
	for i := 0; i < cfg.concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(int64(id) + 1))
			var local []result
			for time.Now().Before(deadline) {
				n := atomic.AddInt64(&counter, 1)
				local = append(local, fire(client, cfg, s, rng, n))
			}
			perWorker[id] = local
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)

	var all []result
	for _, w := range perWorker {
		all = append(all, w...)
	}
	report(cfg, all, elapsed)
}

// fire issues one request based on the scenario and returns its result.
func fire(client *http.Client, cfg config, s seed, rng *rand.Rand, n int64) result {
	switch cfg.scenario {
	case "read":
		return readReq(client, cfg, s, rng)
	case "usage":
		return usageReq(client, cfg, s, n)
	case "invoice":
		return invoiceReq(client, cfg, s)
	default: // mixed: 70% read, 25% usage write, 5% invoice
		r := rng.Intn(100)
		switch {
		case r < 70:
			return readReq(client, cfg, s, rng)
		case r < 95:
			return usageReq(client, cfg, s, n)
		default:
			return invoiceReq(client, cfg, s)
		}
	}
}

func readReq(client *http.Client, cfg config, s seed, rng *rand.Rand) result {
	paths := []string{"/api/v1/plans", "/api/v1/customers", "/api/v1/invoices", "/api/v1/subscriptions", "/healthz"}
	return timed(client, "GET", cfg.baseURL+paths[rng.Intn(len(paths))], s.token, nil)
}

func usageReq(client *http.Client, cfg config, s seed, n int64) result {
	body := map[string]any{
		"customer_id": s.customerID, "subscription_id": s.subID,
		"metric": "api_calls", "quantity": 1,
		"idempotency_key": fmt.Sprintf("load-%d-%d", time.Now().UnixNano(), n),
	}
	return timed(client, "POST", cfg.baseURL+"/api/v1/usage", s.token, body)
}

func invoiceReq(client *http.Client, cfg config, s seed) result {
	body := map[string]any{"subscription_id": s.subID}
	return timed(client, "POST", cfg.baseURL+"/api/v1/invoices/generate", s.token, body)
}

func timed(client *http.Client, method, url, token string, body any) result {
	t0 := time.Now()
	_, code, err := doJSON(client, method, url, token, body)
	res := result{latency: time.Since(t0), status: code}
	if err != nil || code >= 400 {
		res.err = true
	}
	return res
}

func doJSON(client *http.Client, method, url, token string, body any) ([]byte, int, error) {
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, url, rdr)
	if err != nil {
		return nil, 0, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return data, resp.StatusCode, nil
}

func idOf(b []byte) string {
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if v, ok := m["id"].(string); ok {
		return v
	}
	return ""
}

func report(cfg config, all []result, elapsed time.Duration) {
	total := len(all)
	if total == 0 {
		fmt.Println("no requests issued")
		return
	}
	var errs int
	lat := make([]time.Duration, 0, total)
	statusClass := map[string]int{}
	for _, r := range all {
		if r.err {
			errs++
		}
		lat = append(lat, r.latency)
		statusClass[fmt.Sprintf("%dxx", r.status/100)]++
	}
	sort.Slice(lat, func(i, j int) bool { return lat[i] < lat[j] })

	pct := func(p float64) time.Duration {
		idx := int(p / 100 * float64(total))
		if idx >= total {
			idx = total - 1
		}
		return lat[idx]
	}
	var sum time.Duration
	for _, d := range lat {
		sum += d
	}
	avg := sum / time.Duration(total)
	rps := float64(total) / elapsed.Seconds()
	errRate := float64(errs) / float64(total) * 100

	fmt.Printf("  requests     : %d in %.1fs\n", total, elapsed.Seconds())
	fmt.Printf("  throughput   : %.0f req/s\n", rps)
	fmt.Printf("  errors       : %d (%.2f%%)\n", errs, errRate)
	fmt.Printf("  status       : %v\n", statusClass)
	fmt.Printf("  latency avg  : %s\n", round(avg))
	fmt.Printf("  latency p50  : %s\n", round(pct(50)))
	fmt.Printf("  latency p90  : %s\n", round(pct(90)))
	fmt.Printf("  latency p95  : %s\n", round(pct(95)))
	fmt.Printf("  latency p99  : %s\n", round(pct(99)))
	fmt.Printf("  latency max  : %s\n", round(lat[total-1]))
	fmt.Printf("CSV,%s,%d,%d,%.0f,%.2f,%s,%s,%s,%s\n",
		cfg.scenario, cfg.concurrency, total, rps, errRate,
		round(pct(50)), round(pct(90)), round(pct(99)), round(lat[total-1]))
}

func round(d time.Duration) time.Duration {
	if d > time.Millisecond {
		return d.Round(100 * time.Microsecond)
	}
	return d.Round(time.Microsecond)
}
