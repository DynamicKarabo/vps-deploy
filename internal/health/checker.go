package health

import (
	"fmt"
	"net/http"
	"time"
)

// Checker polls an HTTP URL until it returns a 2xx status or the timeout is reached.
type Checker struct {
	URL         string
	Timeout     time.Duration
	RetryDelay  time.Duration
	MaxRetries  int
}

// NewChecker creates a new health checker.
// timeoutSec: max total time to wait for health.
// retryDelaySec: seconds between retries.
func NewChecker(url string, timeoutSec, retryDelaySec int) *Checker {
	return &Checker{
		URL:        url,
		Timeout:    time.Duration(timeoutSec) * time.Second,
		RetryDelay: time.Duration(retryDelaySec) * time.Second,
	}
}

// Poll repeatedly checks the health URL until it returns a 2xx status or the
// timeout is exceeded. Returns nil on success, error if unhealthy or timed out.
func (c *Checker) Poll() error {
	deadline := time.Now().Add(c.Timeout)
	client := &http.Client{Timeout: 10 * time.Second}

	for time.Now().Before(deadline) {
		resp, err := client.Get(c.URL)
		if err != nil {
			fmt.Printf("  health check attempt failed: %v (retrying...)\n", err)
			time.Sleep(c.RetryDelay)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}

		fmt.Printf("  health check returned %d (retrying...)\n", resp.StatusCode)
		time.Sleep(c.RetryDelay)
	}

	return fmt.Errorf("health check timed out after %s", c.Timeout)
}
