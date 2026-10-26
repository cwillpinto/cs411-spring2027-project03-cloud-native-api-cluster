package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func handler(counter Counter, limiter *Limiter) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := counter.Ping(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("/api/work", func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow(time.Now()) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		iterations := 100000
		if value := r.URL.Query().Get("iterations"); value != "" {
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed < 1 || parsed > 2000000 {
				http.Error(w, "invalid iterations", http.StatusBadRequest)
				return
			}
			iterations = parsed
		}
		var checksum uint64
		for i := 0; i < iterations; i++ {
			checksum = checksum*1664525 + uint64(i) + 1013904223
		}
		count, err := counter.Increment(r.Context())
		if err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"requests": count, "checksum": checksum})
	})
	return mux
}
