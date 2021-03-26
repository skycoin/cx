package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type TesterConfig struct {
	cxPath         string
	workingDir     string
	testsMask      Bits
	logMask        Bits
	defaultTimeout time.Duration
}

type tester struct {
	cfg         *TesterConfig
	testCount   int
	testSuccess int
	testSkipped int
}

func NewTester(cfg *TesterConfig) *tester {
	return &tester{
		cfg: cfg,
	}
}

type RunConfig struct {
	Args     string
	ExitCode int
	Desc     string
	Output   string
	Filter   Bits
	Timeout  time.Duration
}

func (t *tester) Run(args string, exitCode int, desc string) {
	t.RunWithConfig(&RunConfig{
		Args:     args,
		ExitCode: exitCode,
		Desc:     desc,
		Output:   "",
		Filter:   TestStable,
		Timeout:  t.cfg.defaultTimeout,
	})
}

func (t *tester) RunEx(args string, exitCode int, desc string, filter Bits, timeout time.Duration) {
	t.RunWithConfig(&RunConfig{
		Args:     args,
		ExitCode: exitCode,
		Desc:     desc,
		Output:   "",
		Filter:   filter,
		Timeout:  timeout,
	})
}

func (t *tester) RunWithConfig(cfg *RunConfig) {
	if cfg.Timeout == 0 {
		cfg.Timeout = t.cfg.defaultTimeout
	}

	if cfg.Filter == TestNone {
		cfg.Filter = TestStable
	}

	if !Has(t.cfg.testsMask, cfg.Filter) {
		if Has(t.cfg.logMask, LogSkip) {
			fmt.Printf("#--- | SKIPPED | na | '%s' | na | na | %s\\n", cfg.Args, cfg.Desc)
		}

		t.testSkipped = t.testSkipped + 1
		return
	}

	cmd := exec.Command(t.cfg.cxPath, strings.Split(cfg.Args, " ")...)
	cmd.Dir = t.cfg.workingDir

	start := time.Now().Unix()
	out, err := runCmd(cmd, cfg.Timeout)
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
				padding(t), t.testCount, timing, cfg.Args)
			return
		}

		if exitError, ok := err.(*exec.ExitError); ok {
			ec = exitError.ExitCode()
			stderr = exitError.Stderr
		}

		if ec != cfg.ExitCode {
			if Has(t.cfg.logMask, LogFail) {
				fmt.Printf("#%s%d | FAILED  | %s | '%s' | exec.Command exited with code %d, expected %d\n",
					padding(t), t.testCount, timing, cfg.Args, ec, cfg.ExitCode)
			}

			if Has(t.cfg.logMask, LogStderr) {
				fmt.Printf("#%s%d | Stderr: %v, %v\n",
					padding(t), t.testCount, string(out), string(stderr))
			}
			return
		}
	}

	if cfg.Output != "" && cfg.Output != string(out) {
		fmt.Printf("#%s%d | FAILED  | %s | '%s' | Got output '%s', expected '%s'\n",
			padding(t), t.testCount, timing, cfg.Args, string(out), cfg.Output)
		return
	}

	if Has(t.cfg.logMask, LogSuccess) {
		fmt.Printf("#%s%d | SUCCESS | %s | '%s' | expected %d | got %d \n",
			padding(t), t.testCount, timing, cfg.Args, cfg.ExitCode, ec)
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
