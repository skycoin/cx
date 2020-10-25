package cxutil

import (
	"flag"
	"fmt"
)

// CommandFunc represents a function that contains a command's logic.
type CommandFunc func(args []string)

// CommandMap maps a command with it's subcommands.
type CommandMap struct {
	cmd *flag.FlagSet
	sig UsageFormat
	m   map[string]CommandFunc
	l   []string
}

func NewCommandMap(cmd *flag.FlagSet, cmdCap int, usageSig UsageFormat) *CommandMap {
	if usageSig == nil {
		usageSig = DefaultUsageFormat("args")
	}

	cm := &CommandMap{
		cmd: cmd,
		sig: usageSig,
		m:   make(map[string]CommandFunc, cmdCap),
		l:   make([]string, 0, cmdCap),
	}

	return cm
}

func (cm *CommandMap) AddSubcommand(subcommand string, fn func(args []string)) *CommandMap {
	cm.m[subcommand] = fn
	cm.l = append(cm.l, subcommand)
	return cm
}

func (cm *CommandMap) ParseAndRun(args []string) int {
	cm.cmd.Usage = func() { cm.sig(cm.cmd, cm.l) }
	if err := cm.cmd.Parse(args); err != nil {
		return 1
	}

	// replace args
	args = cm.cmd.Args()

	if len(args) < 1 {
		CmdErrorf(cm.cmd, fmt.Errorf("command '%s' expects an additional subcommand", cm.cmd.Name()))
		cm.cmd.Usage()
		return 1
	}

	cmdFn, ok := cm.m[args[0]]
	if !ok {
		CmdErrorf(cm.cmd, fmt.Errorf("subcommand '%s' does not exist", args[0]))
		cm.cmd.Usage()
		return 1
	}

	cmdFn(args[1:])
	return 0
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func CmdPrintf(cmd *flag.FlagSet, format string, a ...interface{}) {
	if lastI := len(format)-1; lastI >= 0 && format[lastI] != '\n' {
		format += "\n"
	}
	_, _ = fmt.Fprintf(cmd.Output(), format, a...)
}

func CmdErrorf(cmd *flag.FlagSet, err error) {
	CmdPrintf(cmd, "Error:\n  %v\n", err)
}

func CountDefinedFlags(cmd *flag.FlagSet) (n int) {
	cmd.VisitAll(func(f *flag.Flag) { n++ })
	return n
}
