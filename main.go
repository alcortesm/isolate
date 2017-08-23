package main

import (
	"fmt"
	"os"
	"os/exec"
)

const usage = `Runs the specified command (and its arguments) with configurable
levels of isolation.

The standard input, output and error of the isolate command
are redirected to the provided command.

Usage:
	isolate command [arguments]

The arguments are:
	command:   the command to be run.
	arguments: the command's arguments.
`

func main() {
	cmd, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := run(cmd...); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(0)
}

func parseArgs() (cmd []string, err error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("needs a command")
	}
	return os.Args[1:], nil
}

func run(cmd ...string) error {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
