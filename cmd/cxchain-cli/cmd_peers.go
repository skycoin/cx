package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/SkycoinProject/cx-chains/src/api"

	"github.com/SkycoinProject/cx/cxutil"
)

func cmdPeers(args []string) {
	// spec is the chain spec obtained from ENV
	spec := parseSpecFilepathEnv()

	// rootCmd is the root command of the 'peers' subcommand
	rootCmd := flag.NewFlagSet("cxchain-cli peers", flag.ExitOnError)

	// nodeAddr holds the value parsed from the flags 'node' and 'n'
	nodeAddr := fmt.Sprintf("http://127.0.0.1:%d", spec.Node.WebInterfacePort)
	addNodeAddrFlag := func(cmd *flag.FlagSet) {
		cmd.StringVar(&nodeAddr, "node", nodeAddr, "HTTP API `ADDRESS` of cxchain node")
		cmd.StringVar(&nodeAddr, "n", nodeAddr, "shorthand for 'node'")
	}

	// modCmdPrelude is the prelude logic to 'modify' based commands
	modCmdPrelude := func(name string, args []string, argsName string) *flag.FlagSet {
		cmd := flag.NewFlagSet(name, flag.ExitOnError)
		cmd.Usage = func() {
			cxutil.CmdPrintf(cmd, "Usage:\n  %s [%s...]\n", name, argsName)
			cxutil.CmdPrintf(cmd, "Flags:\n")
			cmd.PrintDefaults()
		}
		addNodeAddrFlag(cmd)

		if len(args) < 1 {
			cxutil.CmdErrorf(cmd, fmt.Errorf("no %s specified", argsName))
			cmd.Usage()
			os.Exit(1)
		}

		if err := cmd.Parse(args); err != nil {
			cxutil.CmdErrorf(cmd, err)
			cmd.Usage()
			os.Exit(1)
		}
		return cmd
	}

	// viewCmdPrelude is the prelude logic to 'view' based commands
	viewCmdPrelude := func(name string, args []string) *flag.FlagSet {
		cmd := flag.NewFlagSet(name, flag.ExitOnError)
		cmd.Usage = func() {
			cxutil.CmdPrintf(cmd, "Usage:\n  %s\n", name)
			cxutil.CmdPrintf(cmd, "Flags:\n")
			cmd.PrintDefaults()
		}
		addNodeAddrFlag(cmd)

		if err := cmd.Parse(args[1:]); err != nil {
			cxutil.CmdErrorf(cmd, err)
			cmd.Usage()
			os.Exit(1)
		}
		return cmd
	}

	// connectionSubCmd contains the 'cxchain-cli peers connection' logic
	connectionSubCmd := func(args []string) {
		cmd := modCmdPrelude("cxchain-cli peers connection", args, "addresses")
		c := api.NewClient(nodeAddr)

		for i, addr := range cmd.Args() {
			log := log.WithField("addr", addr).WithField("i", i)

			out, err := c.NetworkConnection(addr)
			if err != nil {
				log.WithError(err).Error("Connection failed.")
				continue
			}
			j, err := json.MarshalIndent(out, "", "\t")
			if err != nil {
				panic(err)
			}
			log.WithField("data", string(j)).Info("Connected.")
		}
	}

	// disconnectSubCmd contains the 'cxchain-cli peers disconnect' logic
	disconnectSubCmd := func(args []string) {
		cmd := modCmdPrelude("cxchain-cli peers disconnect", args, "conn_ids")
		c := api.NewClient(nodeAddr)

		for i, connIDStr := range cmd.Args() {
			log := log.WithField("conn_id", connIDStr).WithField("i", i)

			connID, err := strconv.ParseUint(connIDStr, 10, 64)
			if err != nil {
				log.WithError(err).Fatal("Failed to parse conn_id.")
			}

			if err := c.Disconnect(connID); err != nil {
				log.WithError(err).Error("Failed to disconnect.")
				continue
			}

			log.Info("Disconnected.")
		}
	}

	// listSubCmd contains the 'cxchain-cli peers list' logic
	listSubCmd := func(args []string) {
		_ = viewCmdPrelude("cxchain-cli peers list", args)
		c := api.NewClient(nodeAddr)

		conns, err := c.NetworkConnections(nil)
		if err != nil {
			log.WithError(err).Fatal("Failed to obtain connections.")
		}

		j, err := json.MarshalIndent(conns, "", "\t")
		if err != nil {
			panic(err)
		}
		log.WithField("data", string(j)).Info()
	}

	// cmdMap contains the map of subcommands.
	cmdMap := cxutil.NewCommandMap(rootCmd, 3, cxutil.DefaultUsageFormat("args")).
		AddSubcommand("connection", connectionSubCmd).
		AddSubcommand("disconnect", disconnectSubCmd).
		AddSubcommand("list", listSubCmd)

	// // Run and exit.
	// if len(args) < 1 {
	// 	cxutil.CmdErrorf(rootCmd, fmt.Errorf("command '%s' expects an additonal subcommand", rootCmd.Name()))
	// 	rootCmd.Usage()
	// }

	os.Exit(cmdMap.ParseAndRun(args))
}
