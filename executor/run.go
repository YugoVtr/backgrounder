//go:generate mockgen -source run.go -destination ../mock/executor.go -package mock
package executor

import (
	"io"
	osexec "os/exec"
)

type Exec interface {
	Run(io.Writer, string, ...string) error
}

func New() exec {
	return exec{}
}

type exec struct{}

func (e exec) Run(out io.Writer, command string, args ...string) error {
	cmd := osexec.Command(command, args...)
	cmd.Stdout = out

	return cmd.Run()
}
