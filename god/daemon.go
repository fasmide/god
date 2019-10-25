package god

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

// Daemon holds entire configuration
type Daemon struct {
	Processes []*Process

	// shuttingDown is set when we are expecting a clean shutdown
	// must not be read or written without the lock
	sync.Mutex
	shuttingDown bool
}

// Run starts all processes and should never stop
func (d *Daemon) Run() error {
	d.handleSignals()

	var group errgroup.Group
	for _, p := range d.Processes {
		p := p
		group.Go(func() error {
			// run the process in question
			err := p.Run()

			// if error is nil - the process exited with code 0
			// this is not usually an error but for us it is
			if err == nil {
				err = errors.New("exit code 0")
			}

			// if we are "in shutdown mode" this is to be expected
			// the exitcode may not be zero - but at least it exited
			d.Lock()
			if d.shuttingDown {
				fmt.Printf("god: shutdown active: %s exited: %s\n", p.Name, err)
				d.Unlock()
				return nil
			}
			d.Unlock()

			fmt.Printf("god: %s exited: %s - initiating shutdown\n", p.Name, err)
			go d.Shutdown()

			// if we end here - the process have exited one way or another
			return fmt.Errorf("%s exited: %s", p.Name, err)
		})
	}

	// Block until all processes have finished
	return group.Wait()
}

// Shutdown sends respective signals to processes
func (d *Daemon) Shutdown() {
	d.Lock()
	d.shuttingDown = true
	d.Unlock()
	for _, p := range d.Processes {
		p.Shutdown()
	}
}

// handleSignals listens for SIGTERM and runs .Shutdown on all processes
func (d *Daemon) handleSignals() {
	go func() {
		termSignal := make(chan os.Signal, 1)
		signal.Notify(termSignal, syscall.SIGTERM)

		// Block until a signal is received.
		s := <-termSignal
		fmt.Println("god: shutting down, received signal:", s)
		d.Shutdown()
	}()
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
