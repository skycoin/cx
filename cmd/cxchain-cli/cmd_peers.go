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
	rootCmd := flag.NewFlagSet(args[0], flag.ExitOnError)

	// nodeAddr holds the value parsed from the flags 'node' and 'n'
	nodeAddr := fmt.Sprintf("http://127.0.0.1:%d", spec.Node.WebInterfacePort)
	addNodeAddrFlag := func(cmd *flag.FlagSet) {
		cmd.StringVar(&nodeAddr, "node", nodeAddr, "HTTP API `ADDRESS` of cxchain node")
		cmd.StringVar(&nodeAddr, "n", nodeAddr, "shorthand for 'node'")
	}

	// modCmdPrelude is the prelude logic to 'modify' based commands
	modCmdPrelude := func(args []string, argsName string) *flag.FlagSet {
		cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
		cmd.Usage = func() {
			cxutil.CmdPrintf(cmd, "Usage:\n")
			cxutil.CmdPrintf(cmd, "  %s %s [%s...]\n", os.Args[0], os.Args[1], argsName)
			cxutil.CmdPrintf(cmd, "Flags:\n")
			cmd.PrintDefaults()
		}
		addNodeAddrFlag(cmd)

		if len(args) < 2 {
			cxutil.CmdPrintf(cmd, "Error:\n")
			cxutil.CmdPrintf(cmd, "  no %s specified\n", argsName)
			cmd.Usage()
			os.Exit(1)
		}

		if err := cmd.Parse(args[1:]); err != nil {
			cxutil.CmdPrintf(cmd, "Error:\n")
			cxutil.CmdPrintf(cmd, "  %v\n", err)
			cmd.Usage()
			os.Exit(1)
		}
		return cmd
	}

	// viewCmdPrelude is the prelude logic to 'view' based commands
	viewCmdPrelude := func(args []string) *flag.FlagSet {
		cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
		cmd.Usage = func() {
			cxutil.CmdPrintf(cmd, "Usage:\n")
			cxutil.CmdPrintf(cmd, "  %s %s\n", os.Args[0], os.Args[1])
			cxutil.CmdPrintf(cmd, "Flags:\n")
			cmd.PrintDefaults()
		}
		addNodeAddrFlag(cmd)

		if err := cmd.Parse(args[1:]); err != nil {
			cxutil.CmdPrintf(cmd, "Error:\n")
			cxutil.CmdPrintf(cmd, "  %v\n", err)
			cmd.Usage()
			os.Exit(1)
		}
		return cmd
	}

	// connectSubCmd contains the 'cxchain-cli peers connect' logic
	connectSubCmd := func(args []string) {
		cmd := modCmdPrelude(args, "addresses")
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
		cmd := modCmdPrelude(args, "conn_ids")
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
		_ = viewCmdPrelude(args)
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
	cmdMap := cxutil.NewCommandMap(rootCmd, 3, cxutil.DefaultUsageSignature()).
		Add("connect", connectSubCmd).
		Add("disconnect", disconnectSubCmd).
		Add("list", listSubCmd)

	// Run and exit.
	os.Exit(cmdMap.ParseAndRun(args[1:]))
}
