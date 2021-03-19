package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	cxPath         string
	workingDir     string
	testsMask      Bits
	logMask        Bits
	defaultTimeout time.Duration
}

type tester struct {
	cfg         *Config
	testCount   int
	testSuccess int
	testSkipped int
}

func NewTester(cfg *Config) *tester {
	return &tester{
		cfg: cfg,
	}
}

func (t *tester) Run(args string, exitCode int, desc string) {
	t.RunEx(args, exitCode, desc, TestStable, t.cfg.defaultTimeout)
}

func (t *tester) RunEx(args string, exitCode int, desc string, filter Bits, timeout time.Duration) {
	if timeout == 0 {
		timeout = t.cfg.defaultTimeout
	}

	if !Has(t.cfg.testsMask, filter) {
		if Has(t.cfg.logMask, LogSkip) {
			fmt.Printf("#--- | SKIPPED | na | '%s' | na | na | %s\\n", args, desc)
		}

		t.testSkipped = t.testSkipped + 1
		return
	}

	cmd := exec.Command(t.cfg.cxPath, strings.Split(args, " ")...)
	cmd.Dir = t.cfg.workingDir

	start := time.Now().Unix()
	out, err := runCmd(cmd, timeout)
	end := time.Now().Unix()

	timing := "na"
	if Has(t.cfg.logMask, LogTime) {
		timing = fmt.Sprintf("%dms", end-start)
	}

	t.testCount += 1

	var ec int
	var stderr []byte
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Printf("#%s%d | FAILED  | %s | '%s' | exec.Command timeout\n",
				padding(t), t.testCount, timing, args)
			return
		}

		if exitError, ok := err.(*exec.ExitError); ok {
			ec = exitError.ExitCode()
			stderr = exitError.Stderr
		}

		if ec != exitCode {
			if Has(t.cfg.logMask, LogFail) {
				fmt.Printf("#%s%d | FAILED  | %s | '%s' | exec.Command exited with code %d, expected %d\n",
					padding(t), t.testCount, timing, args, ec, exitCode)
			}

			if Has(t.cfg.logMask, LogStderr) {
				fmt.Printf("#%s%d | Stderr: %v, %v\n",
					padding(t), t.testCount, string(out), string(stderr))
			}
			return
		}
	}

	if Has(t.cfg.logMask, LogSuccess) {
		fmt.Printf("#%s%d | SUCCESS | %s | '%s' | expected %d | got %d \n",
			padding(t), t.testCount, timing, args, exitCode, ec)
	}
	t.testSuccess += 1
}

func runCmd(cmd *exec.Cmd, timeout time.Duration) ([]byte, error) {
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// This time we can simply use Output() to get the result.
	out, err := cmd.Output()

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		return out, ctx.Err()
	}

	// If there's no context error, we know the command completed (or errored).
	return out, err
}

func padding(t *tester) string {
	var padding string
	if t.testCount < 10 {
		padding = "  "
	} else if t.testCount < 100 {
		padding = " "
	}
	return padding
}
