package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	app := &cli.App{
		Name:  "cxtest",
		Usage: "cx programs tester",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cxpath",
				Usage:       "cx binary path",
				DefaultText: "../bin/cx",
			},
			&cli.StringFlag{
				Name:        "wdir",
				Usage:       "working directory with *.cx tests",
				DefaultText: "./tests",
			},
			&cli.StringFlag{
				Name:  "log",
				Usage: "Enable logMask set (all, success, stderr, fail, skip, time)",
			},
			&cli.StringFlag{
				Name:  "enable-tests",
				Usage: "Enable test set (all, stable, issue, gui)",
			},
			&cli.StringFlag{
				Name:  "disable-tests",
				Usage: "Disable test set (all, stable, issue, gui)",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Print debug information",
			},
		},
		Action: func(c *cli.Context) error {
			return Execute(c)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Execute(c *cli.Context) error {
	cxPath := c.String("cxpath")
	if cxPath == "" {
		cxPath = "../bin/cx"
	}

	workingDir := c.String("wdir")
	if workingDir == "" {
		workingDir = "./tests"
	}

	debug := c.Bool("debug")

	var parseBitMask = func(flagName string, bitsMap map[string]Bits, defaultBit Bits) Bits {
		var mask Bits = 0
		flags := strings.Split(c.String(flagName), ",")
		for _, flag := range flags {
			mask = Set(mask, bitsMap[flag])
		}
		if debug {
			fmt.Printf("Parsed bit mask for flag %s%s: %06b\n", flagName, flags, mask)
		}
		return mask
	}

	logMask := parseBitMask("log", LogBits, LogNone)
	enableTestsMaks := parseBitMask("enable-tests", TestBits, TestNone)
	disabledTestsMaks := parseBitMask("disable-tests", TestBits, TestNone)

	if enableTestsMaks == TestAll && disabledTestsMaks == TestAll {
		return errors.New("invalid test flags combination")
	}

	var testsMask Bits = TestAll
	if debug {
		fmt.Printf("Initial test mask: %06b\n", testsMask)
	}
	// turn on only enabled tests
	if enableTestsMaks != TestNone {
		testsMask = testsMask & enableTestsMaks
	}
	// turn off only disabled test
	testsMask = testsMask &^ disabledTestsMaks
	if debug {
		fmt.Printf("Resulting test mask: %06b\n", testsMask)
	}

	tester := NewTester(&Config{
		cxPath:         cxPath,
		workingDir:     workingDir,
		testsMask:      testsMask,
		logMask:        logMask,
		defaultTimeout: 10 * time.Second,
	})

	var start = time.Now().Unix()

	fmt.Printf("Running CX tests in dir: '%s'\n", workingDir)
	runTestCases(tester)
	end := time.Now().Unix()

	if Has(logMask, LogTime) {
		fmt.Printf("\nTests finished after %d milliseconds", end-start)
	}

	fmt.Printf("\nA total of %d tests were performed\n", tester.testCount)
	fmt.Printf("%d were successful\n", tester.testSuccess)
	fmt.Printf("%d failed\n", tester.testCount-tester.testSuccess)
	fmt.Printf("%d skipped\n", tester.testSkipped)

	if tester.testCount == 0 || (tester.testSuccess != tester.testCount) {
		return errors.New("not all test succeeded")
	}

	return nil
}
