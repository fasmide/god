package god

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Process represents a single process we would like to run
type Process struct {
	Name         string
	Cmd          string
	Bash         bool
	Requirements Requirements `yaml:"requires"`
}

// Run execs the command and blocks until it have exited
func (p *Process) Run() error {
	// wait for requirements to fulfill
	err := p.Requirements.Wait()
	if err != nil {
		return fmt.Errorf("requirements failed for %s: %s", p.Name, err)
	}

	var cmd *exec.Cmd

	if p.Bash {
		cmd = exec.Command("bash", "-c", p.Cmd)
	} else {
		fields := strings.Fields(p.Cmd)
		cmd = exec.Command(fields[0], fields[1:]...)
	}

	relay := func(r io.Reader, w io.Writer) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			fmt.Fprintf(w, "%s: %s\n", p.Name, scanner.Text())

		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("god: cannot open stream for scanning %s: %s\n", p.Name, err)
		}

	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go relay(stderr, os.Stderr)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go relay(stdout, os.Stdout)

	return cmd.Run()
}
