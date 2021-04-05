package runner

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	CxPath         string
	WorkingDir     string
	TestsMask      Bits
	LogMask        Bits
	DefaultTimeout time.Duration
}

type TestRunner struct {
	cfg         *Config
	TestCount   int
	TestSuccess int
	TestSkipped int
}

func NewTestRunner(cfg *Config) *TestRunner {
	return &TestRunner{
		cfg: cfg,
	}
}

func (t *TestRunner) Run(args string, exitCode int, desc string) {
	t.RunEx(args, exitCode, desc, TestStable, t.cfg.DefaultTimeout)
}

func (t *TestRunner) RunEx(args string, exitCode int, desc string, filter Bits, timeout time.Duration) {
	if timeout == 0 {
		timeout = t.cfg.DefaultTimeout
	}

	if !Has(t.cfg.TestsMask, filter) {
		if Has(t.cfg.LogMask, LogSkip) {
			fmt.Printf("#--- | SKIPPED | na | '%s' | na | na | %s\\n", args, desc)
		}

		t.TestSkipped = t.TestSkipped + 1
		return
	}

	cmd := exec.Command(t.cfg.CxPath, strings.Split(args, " ")...)
	cmd.Dir = t.cfg.WorkingDir

	start := time.Now().Unix()
	out, err := runCmd(cmd, timeout)
	end := time.Now().Unix()

	timing := "na"
	if Has(t.cfg.LogMask, LogTime) {
		timing = fmt.Sprintf("%dms", end-start)
	}

	t.TestCount += 1

	var ec int
	var stderr []byte
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Printf("#%s%d | FAILED  | %s | '%s' | exec.Command timeout\n",
				padding(t), t.TestCount, timing, args)
			return
		}

		if exitError, ok := err.(*exec.ExitError); ok {
			ec = exitError.ExitCode()
			stderr = exitError.Stderr
		}

		if ec != exitCode {
			if Has(t.cfg.LogMask, LogFail) {
				fmt.Printf("#%s%d | FAILED  | %s | '%s' | exec.Command exited with code %d, expected %d\n",
					padding(t), t.TestCount, timing, args, ec, exitCode)
			}

			if Has(t.cfg.LogMask, LogStderr) {
				fmt.Printf("#%s%d | Stderr: %v, %v\n",
					padding(t), t.TestCount, string(out), string(stderr))
			}
			return
		}
	}

	if Has(t.cfg.LogMask, LogSuccess) {
		fmt.Printf("#%s%d | SUCCESS | %s | '%s' | expected %d | got %d \n",
			padding(t), t.TestCount, timing, args, exitCode, ec)
	}
	t.TestSuccess += 1
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

func padding(t *TestRunner) string {
	var padding string
	if t.TestCount < 10 {
		padding = "  "
	} else if t.TestCount < 100 {
		padding = " "
	}
	return padding
}
