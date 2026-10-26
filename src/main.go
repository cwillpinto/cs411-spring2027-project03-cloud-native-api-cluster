package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	redis := os.Getenv("REDIS_ADDR")
	if redis == "" {
		redis = "redis:6379"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, handler(RedisCounter{Address: redis}, NewLimiter(20, time.Second))))
}
