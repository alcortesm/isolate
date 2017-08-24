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

The standard input, output and error of the isolated command
are redirected to the ones of the main "isolate" process.

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
		"prints the exit code of the isolated command to stdout.")
	opts := runOpts{
		newUTS: flag.Bool("uts", false,
			"run the isolated command in a new UTS namespace, initialized after the one of the main process."),
	}
	flag.Parse()

	exitCode, err := run(flag.Args(), opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *showExitCode {
		fmt.Println("Exit Code", exitCode)
	}
}

type runOpts struct {
	newUTS *bool
}

func run(words []string, opts runOpts) (int, error) {
	if len(words) == 0 {
		return 0, fmt.Errorf("missing command")
	}
	cmd := exec.Command(words[0], words[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = new(syscall.SysProcAttr)
	if *opts.newUTS {
		// requires Linux >= 2.6.19 and CAP_SYS_ADMIN
		cmd.SysProcAttr.Cloneflags = syscall.CLONE_NEWUTS
	}
	return unixExitCodeOrError(cmd.Run())
}

func unixExitCodeOrError(err error) (int, error) {
	var exitErr *exec.ExitError
	var ok bool
	if exitErr, ok = err.(*exec.ExitError); !ok {
		return 0, err
	}
	var status syscall.WaitStatus
	if status, ok = exitErr.Sys().(syscall.WaitStatus); !ok {
		return 0, fmt.Errorf(
			"unsupported (non Unix) system-dependent exit information")
	}
	return status.ExitStatus(), nil
}
