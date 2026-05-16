package health

import (
	"context"
	"fmt"
	"net"
	"time"

	"ClusterGuard/internal/domain"
)

// TCPChecker probes reachability via TCP dial.
type TCPChecker struct {
	Timeout time.Duration
}

func NewTCPChecker(timeout time.Duration) *TCPChecker {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &TCPChecker{Timeout: timeout}
}

func (c *TCPChecker) Type() string { return "tcp" }

func (c *TCPChecker) Check(ctx context.Context, host string, port int, _ string) (bool, int64, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	dialer := net.Dialer{Timeout: c.Timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, time.Since(start).Milliseconds(), nil
	}
	_ = conn.Close()
	return true, time.Since(start).Milliseconds(), nil
}

var _ domain.HealthChecker = (*TCPChecker)(nil)
