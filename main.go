package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const usage = `Runs the specified command (and its arguments) with configurable
levels of isolation.

The standard input, output and error of the isolate command
are redirected to the provided command.

Usage:
	isolate [options] command [arguments]

The arguments are:
  command
	the command to be run.
  arguments
	the command's arguments.

The options are:
`

func main() {
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
	}
	showExitCode := flag.Bool("exitCode", false,
		"prints the exit code of the isolated command to stdout")
	flag.Parse()
	cmd := flag.Args()
	if len(cmd) == 0 {
		fmt.Fprintln(os.Stderr, "needs a command")
		os.Exit(1)
	}
	exitCode, err := run(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *showExitCode {
		fmt.Println("Exit Code", exitCode)
	}
	os.Exit(0)
}

func run(words []string) (int, error) {
	cmd := exec.Command(words[0], words[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return unixExitCodeOrError(cmd.Run())
}

func unixExitCodeOrError(err error) (int, error) {
	var exitErr *exec.ExitError
	var ok bool
	if exitErr, ok = err.(*exec.ExitError); !ok {
		return 0, err
	}
	var status syscall.WaitStatus // Unix
	if status, ok = exitErr.Sys().(syscall.WaitStatus); !ok {
		return 0, fmt.Errorf(
			"unsupported (non Unix) system-dependent exit information")
	}
	return status.ExitStatus(), nil
}
