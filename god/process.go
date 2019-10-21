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
	Name string
	Cmd  string
	Bash bool
}

// Run execs the command and blocks until it have exited
func (p *Process) Run() error {
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
			panic(fmt.Sprintf("could not scan %s: %s", p.Name, err))
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
			panic(fmt.Sprintf("could not scan %s: %s", p.Name, err))
		}
	}()

	return cmd.Run()
}
