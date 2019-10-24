package god

import "golang.org/x/sync/errgroup"

// Requirements is a slice of requires
type Requirements []Requires

// Wait waits for all requirements to be fulfilled
func (r Requirements) Wait() error {
	var errGroup errgroup.Group

	for _, q := range r {
		errGroup.Go(q.Wait)
	}

	return errGroup.Wait()

}
