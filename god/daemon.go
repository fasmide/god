package god

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Daemon holds entire configuration
type Daemon struct {
	Processes []Process
}

// Run starts all processes and should never stop
func (d *Daemon) Run() error {
	failChannel := make(chan error)

	for _, p := range d.Processes {
		go func(p Process) {
			err := p.Run()

			// if we end here - the process have exited one way or another
			failChannel <- fmt.Errorf("%s exited: %s", p.Name, err)
		}(p)
	}

	// Block until a process stops
	err := <-failChannel
	return err
}

// Load loads yaml and returns a daemon ready to run
func Load(path string) (*Daemon, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("io error: %w", err)
	}

	c := &Daemon{}

	err = yaml.Unmarshal(d, c)
	if err != nil {
		return nil, fmt.Errorf("yml error: %w", err)
	}
	return c, nil
}
