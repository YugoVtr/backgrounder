//go:generate mockgen -source run.go -destination ../mock/executor.go -package mock
package exec

import (
	"os/exec"
)

// Commander is an abstrataction to execute the named program with
// the given arguments.
type Commander func(s string, a ...string) Cmd

// Command is an adapter to os/exec.Command
// returning an abstraction rather than a concrete implementation.
func Command(s string, a ...string) Cmd {
	return exec.Command(s, a...)
}

// Cmd represents an external command being prepared or run.
type Cmd interface {
	Run() error
	Output() ([]byte, error)
}
