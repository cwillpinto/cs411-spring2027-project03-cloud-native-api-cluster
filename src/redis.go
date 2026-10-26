package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Counter interface {
	Increment(context.Context) (int64, error)
	Ping(context.Context) error
}

// RedisCounter uses a minimal RESP INCR call to avoid an SDK dependency.
type RedisCounter struct{ Address string }

func (r RedisCounter) Ping(ctx context.Context) error {
	dialer := net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", r.Address)
	if err != nil {
		return err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err = fmt.Fprint(conn, "*1\r\n$4\r\nPING\r\n"); err != nil {
		return err
	}
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	if strings.TrimSpace(line) != "+PONG" {
		return fmt.Errorf("unexpected Redis reply: %q", line)
	}
	return nil
}

func (r RedisCounter) Increment(ctx context.Context) (int64, error) {
	dialer := net.Dialer{Timeout: 2 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", r.Address)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err = fmt.Fprint(conn, "*2\r\n$4\r\nINCR\r\n$13\r\nrequest_count\r\n"); err != nil {
		return 0, err
	}
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return 0, err
	}
	if !strings.HasPrefix(line, ":") {
		return 0, fmt.Errorf("unexpected Redis reply: %q", line)
	}
	return strconv.ParseInt(strings.TrimSpace(line[1:]), 10, 64)
}
