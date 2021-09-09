//go:generate mockgen -source run.go -destination ../mock/executor.go -package mock
package exec

import (
	"os/exec"
)

type Commander func(s string, a ...string) Cmd

func Command(s string, a ...string) Cmd {
	return exec.Command(s, a...)
}

type Cmd interface {
	Run() error
	Output() ([]byte, error)
}
