package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const usage = `Runs the specified command (and its arguments) with configurable
levels of namespace isolation.

The standard input, output and error of the isolated command are
redirected to the ones of the main process (the one running the
"isolate" binary).

Usage:
	isolate [options] command [arguments]

The arguments are:
  command
	the command to be run.
  arguments
	the command's arguments.

The options are:
` // ... continued by flag.PrintDefaults().

func main() {
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
	}
	showExitCode := flag.Bool("exitCode", false,
		"Prints the exit code of the isolated command to stdout.")
	opts := runOpts{
		newUTS: flag.Bool("uts", false,
			"Run the isolated command in a new UTS namespace, initialized\n"+
				"\tafter the one of the main process. Requires CAP_SYS_ADMIN and\n"+
				"\tLinux >= 2.6.19."),
		chroot: flag.String("dir", "",
			"Change the root directory to the given directory and use it\n"+
				"\tas the working directory. Requires CAP_SYS_CHROOT."),
		userns: flag.Bool("userns", false,
			"Run the isolated command in a new user namespace, with a\n"+
				"\tfull set of capabilities.  No privileges needed if Linux >= 3.8."),
		pid: flag.Bool("pid", false,
			"Run the command in its own process ID numberspace."),
		mount: flag.Bool("mount", false,
			"Run the comand int its own mount namespace."),
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
	chroot *string
	userns *bool
	pid    *bool
	mount  *bool
}

func run(words []string, opts runOpts) (int, error) {
	if len(words) == 0 {
		return 0, fmt.Errorf("missing command")
	}
	cmd := exec.Command(words[0], words[1:]...)
	//cmd := &exec.Cmd{Path: words[0], Args: words}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = new(syscall.SysProcAttr)
	if *opts.newUTS {
		// requires Linux >= 2.6.19 and CAP_SYS_ADMIN
		cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUTS
	}
	if *opts.chroot != "" {
		cmd.Dir = "/"
		// requires CAP_SYS_CHROOT
		cmd.SysProcAttr.Chroot = *opts.chroot
	}
	if *opts.userns {
		cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWUSER
	}
	if *opts.pid {
		cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWPID
	}
	if *opts.mount {
		fmt.Println("newns")
		cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNS
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
