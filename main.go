package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const usage = `usage:

	isolate command [arguments]

Runs a command (and its arguments) with configurable levels of isolation.

	command:   the command to be run.
	arguments: the command's arguments.
`

func main() {
	cmd, err := parseArgs()
	if err != nil {
		log.Println(err)
		fmt.Print(usage)
		os.Exit(1)
	}
	if err := run(cmd...); err != nil {
		log.Fatalln(err)
	}
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
