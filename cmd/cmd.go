package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	args []string
	err  error
}

func NewCommand(args []string) *Command {
	return &Command{args: args}
}

func (c *Command) Error() error {
	return c.err
}

func (c *Command) String() string {
	return strings.Join(c.args, " ")
}

// Run the specified command piping Stdin, Stdout and Stderr to the parent process.
// Closes the quit channel on exit.
func (c *Command) RunCommand(quit chan int) {
	defer close(quit)
	cmd := &exec.Cmd{Path: c.args[0], Args: c.args, Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	if err := cmd.Run(); err != nil {
		log.Println("Command", c.args, "exited with error", err)
		c.err = err
	}
}
