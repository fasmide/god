package god

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Process represents a single process we would like to run
type Process struct {
	Name     string
	Cmd      string
	Bash     bool
	Requires *Requires
}

// Run execs the command and blocks until it have exited
func (p *Process) Run() error {

	// is this process have a requirement
	// wait for it to be fulfilled
	if p.Requires != nil {
		err := p.Requires.Wait()
		if err != nil {
			return fmt.Errorf("%s: %w", p.Name, err)
		}
	}

	var cmd *exec.Cmd

	if p.Bash {
		cmd = exec.Command("bash", "-c", p.Cmd)
	} else {
		fields := strings.Fields(p.Cmd)
		cmd = exec.Command(fields[0], fields[1:]...)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintf(os.Stderr, "%s: %s\n", p.Name, scanner.Text())

		}

		if err := scanner.Err(); err != nil {
			// panic for now
			fmt.Printf("god: cannot open stderr for scanning %s: %s\n", p.Name, err)
		}
	}()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Fprintf(os.Stdout, "%s: %s\n", p.Name, scanner.Text())

		}

		if err := scanner.Err(); err != nil {
			// panic for now
			fmt.Printf("god: cannot open stderr for scanning %s: %s\n", p.Name, err)
		}
	}()

	return cmd.Run()
}
