package cxutil

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

type UsageSignatureFunc func(cmd *flag.FlagSet, l []string)

func DefaultUsageSignature() UsageSignatureFunc {
	return func(cmd *flag.FlagSet, l []string) {
		sort.Strings(l)
		CmdPrintf(cmd, "Usage:\n")
		CmdPrintf(cmd, "  %s [%s] [args...]\n", cmd.Name(), strings.Join(l, "|"))

		if CountDefinedFlags(cmd) > 0 {
			CmdPrintf(cmd, "Flags:\n")
			cmd.PrintDefaults()
		}
	}
}

type CommandFunc func(args []string)

type CommandMap struct {
	cmd  *flag.FlagSet
	sig  UsageSignatureFunc
	m    map[string]CommandFunc
	l    []string
}

func NewCommandMap(cmd *flag.FlagSet, cmdCap int, usageSig UsageSignatureFunc) *CommandMap {
	if usageSig == nil {
		usageSig = DefaultUsageSignature()
	}

	cm := &CommandMap{
		cmd: cmd,
		sig: usageSig,
		m:   make(map[string]CommandFunc, cmdCap),
		l:   make([]string, 0, cmdCap),
	}

	return cm
}

func (cm *CommandMap) Add(subcommand string, fn func(args []string)) *CommandMap {
	cm.m[subcommand] = fn
	cm.l = append(cm.l, subcommand)
	return cm
}

func (cm *CommandMap) ParseAndRun(args []string) int {
	cm.cmd.Usage = func() { cm.sig(cm.cmd, cm.l) }
	if err := cm.cmd.Parse(args); err != nil {
		return 1
	}

	if len(args) < 1 {
		CmdPrintf(cm.cmd, "Error:\n")
		CmdPrintf(cm.cmd, "  no args specified\n")
		cm.cmd.Usage()
		return 1
	}

	cmdFn, ok := cm.m[args[0]]
	if !ok {
		CmdPrintf(cm.cmd, "Error:\n")
		CmdPrintf(cm.cmd, "  subcommand '%s' does not exist\n", args[0])
		cm.cmd.Usage()
		return 1
	}

	cmdFn(cm.cmd.Args())
	return 0
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func CmdPrintf(cmd *flag.FlagSet, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(cmd.Output(), format, a...)
}

func CountDefinedFlags(cmd *flag.FlagSet) (n int) {
	cmd.VisitAll(func(f *flag.Flag) { n++ })
	return n
}
