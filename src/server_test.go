package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeCounter struct{ count int64 }

func (f *fakeCounter) Increment(context.Context) (int64, error) { f.count++; return f.count, nil }
func (f *fakeCounter) Ping(context.Context) error               { return nil }

func TestRateLimitWindow(t *testing.T) {
	l := NewLimiter(2, time.Second)
	start := time.Unix(100, 0)
	if !l.Allow(start) || !l.Allow(start.Add(100*time.Millisecond)) {
		t.Fatal("initial requests rejected")
	}
	if l.Allow(start.Add(200 * time.Millisecond)) {
		t.Fatal("third request should be limited")
	}
	if !l.Allow(start.Add(time.Second)) {
		t.Fatal("new window did not reset")
	}
}

func TestHealthAndWork(t *testing.T) {
	f := &fakeCounter{}
	h := handler(f, NewLimiter(1, time.Hour))
	for _, tc := range []struct {
		path string
		want int
	}{{"/health", 200}, {"/ready", 200}, {"/api/work?iterations=10", 200}, {"/api/work?iterations=10", 429}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if w.Code != tc.want {
			t.Fatalf("%s status=%d want=%d", tc.path, w.Code, tc.want)
		}
	}
	if f.count != 1 {
		t.Fatalf("database increments=%d, want 1", f.count)
	}
}
