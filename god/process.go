package god

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Process represents a single process we would like to run
type Process struct {
	Name         string
	Cmd          string
	Bash         bool
	StopSignal   string       `yaml: "stop_signal"`
	Requirements Requirements `yaml:"requires"`

	command *exec.Cmd
}

// Run execs the command and blocks until it have exited
func (p *Process) Run() error {
	// wait for requirements to fulfill
	err := p.Requirements.Wait()
	if err != nil {
		return fmt.Errorf("requirements failed for %s: %s", p.Name, err)
	}

	if p.Bash {
		p.command = exec.Command("bash", "-c", p.Cmd)
	} else {
		fields := strings.Fields(p.Cmd)
		p.command = exec.Command(fields[0], fields[1:]...)
	}

	relay := func(r io.Reader, w io.Writer) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			fmt.Fprintf(w, "%s: %s\n", p.Name, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("god: %s: cannot open stream for scanning: %s\n", p.Name, err)
		}

	}

	stderr, err := p.command.StderrPipe()
	if err != nil {
		return err
	}
	go relay(stderr, os.Stderr)

	stdout, err := p.command.StdoutPipe()
	if err != nil {
		return err
	}
	go relay(stdout, os.Stdout)

	return p.command.Run()
}

// Shutdown sends SIGQUIT or configured signal to the process
func (p *Process) Shutdown() {
	signal, exists := signalMap[p.StopSignal]
	if !exists {
		fmt.Printf("god: %s wants to stop with signal %s - there is no such thing - sending SIGQUIT instead\n", p.Name, p.StopSignal)
		signal = syscall.SIGQUIT
	}

	if p.command == nil {
		fmt.Printf("god: should send signal %s to %s but it have not been started\n", signal, p.Name)
		return
	}

	err := p.command.Process.Signal(signal)
	if err != nil {
		fmt.Printf("god: failed to send %s to %s: %s\n", signal, p.Name, err)
	}
}
