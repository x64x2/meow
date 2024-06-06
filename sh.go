package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"codeberg.org/gruf/go-byteutil"
)

// environ is our own copy of the executing environment.
var environ = syscall.Environ()

// shell contains the determined shell binary to use.
var shell = func() string {
	for _, kv := range environ {
		if strings.HasPrefix(kv, "SHELL=") {
			return kv[6:]
		}
	}
	return "/bin/sh"
}()

// dir contains the determined current working directory.
var dir = func() string {
	for _, kv := range environ {
		if strings.HasPrefix(kv, "PWD=") {
			return kv[4:]
		}
	}
	wd, _ := syscall.Getenv("PWD")
	return wd
}()

// cmd crafts a new exec.Cmd from shell expression and environment.
func cmd(expr string, env ...string) *exec.Cmd {
	cmd := new(exec.Cmd)
	cmd.Path = shell
	cmd.Args = []string{shell, "-c", expr}
	cmd.Dir = dir
	cmd.Env = append(cmd.Env, environ...)
	cmd.Env = append(cmd.Env, env...)
	return cmd
}

type ShellExpr string

func (expr ShellExpr) Match(in string) (bool, error) {
	cmd := cmd(string(expr))

	// Prepare standard file bufs.
	errbuf := new(byteutil.Buffer)
	inputStr := strings.NewReader(in)
	cmd.Stdin = inputStr
	cmd.Stdout = io.Discard
	cmd.Stderr = errbuf

	// Run command.
	err := cmd.Run()

	if s := cmd.ProcessState; s == nil {
		// Nil process state => error calling cmd.
		return false, fmt.Errorf("error executing %s: %w", expr, err)

	} else if n := s.ExitCode(); n != 0 && errbuf.Len() > 0 {
		// Non-zero exit code and non-empty stderr => cmd error.
		return false, fmt.Errorf("returned %d: %s", n, bufferOutStr(errbuf))

	} else {
		// Check for match.
		return (n == 0), nil
	}
}

func (expr ShellExpr) Output(env []string, in io.Reader) (string, error) {
	cmd := cmd(string(expr), env...)

	// Prepare standard file bufs.
	outbuf := new(byteutil.Buffer)
	errbuf := new(byteutil.Buffer)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf
	cmd.Stdin = in

	// Run command.
	err := cmd.Run()

	if s := cmd.ProcessState; s == nil {
		// Nil process state => error calling cmd.
		return "", fmt.Errorf("error executing %s: %w", expr, err)

	} else if n := s.ExitCode(); n != 0 && errbuf.Len() > 0 {
		// Non-zero exit code and non-empty stderr => cmd error.
		return "", fmt.Errorf("returned %d: %s", n, bufferOutStr(errbuf))

	}

	// Return stdout always.
	return bufferOutStr(outbuf), nil
}

func bufferOutStr(buf *byteutil.Buffer) string {
	return strings.TrimSpace(buf.String())
}

func (expr *ShellExpr) Set(in string) error {
	*expr = ShellExpr(in)
	return nil
}

func (expr *ShellExpr) Kind() string {
	return "shell-expr"
}

func (expr *ShellExpr) String() string {
	return string(*expr)
}
