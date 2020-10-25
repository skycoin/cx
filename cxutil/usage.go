package cxutil

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

// UsageFormat represents a function that prints command usage.
type UsageFormat func(cmd *flag.FlagSet, subcommands []string)

// DefaultUsageFormat is the default UsageFormat.
func DefaultUsageFormat(argsStr ...string) UsageFormat {
	return func(cmd *flag.FlagSet, subcommands []string) {
		// print: Usage
		PrintCmdUsage(cmd, "Usage", subcommands, argsStr)

		// print: Flags
		PrintCmdFlags(cmd, "Flags")
	}
}

// PrintCmdUsage prints command usage.
func PrintCmdUsage(cmd *flag.FlagSet, heading string, subcommands, flags []string) {
	usageArr := make([]string, 1, 3)
	usageArr[0] = cmd.Name()

	// pre-process: determine subcommands string
	if len(subcommands) > 0 {
		sort.Strings(subcommands)
		usageArr = append(usageArr, fmt.Sprintf("[%s]", strings.Join(subcommands, "|")))
	}

	// pre-process: determine flags string
	for _, s := range flags {
		usageArr = append(usageArr, fmt.Sprintf("[%s...]", s))
	}

	// print: Usage
	CmdPrintf(cmd, "%s:\n  %s\n", heading, strings.Join(usageArr, " "))
}

func PrintCmdFlags(cmd *flag.FlagSet, heading string) {
	if CountDefinedFlags(cmd) > 0 {
		CmdPrintf(cmd, "%s:\n", heading)
		cmd.PrintDefaults()
	}
}