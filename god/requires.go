package god

import (
	"fmt"
	"os"
	"time"
)

// Requires represents (for now) a requirement for a
// filesystem resource
type Requires struct {
	Path    string `yaml:"exists"`
	Timeout time.Duration
}

// Wait blocks until requirements are fulfilled
func (r *Requires) Wait() error {
	start := time.Now()
	for {
		if time.Since(start) > r.Timeout {
			return fmt.Errorf("timed out waiting for %s: %s", r.Path, r.Timeout)
		}

		_, err := os.Stat(r.Path)

		// is this files does not exist - wait and try again
		if os.IsNotExist(err) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		// fail on any other error
		if err != nil {
			return fmt.Errorf("unable to wait for %s: %w", r.Path, err)
		}

		return nil
	}
}
