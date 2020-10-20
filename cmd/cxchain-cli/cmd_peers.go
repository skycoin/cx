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
	spec := parseSpecFilepathEnv()

	rootCmd := flag.NewFlagSet(args[0], flag.ExitOnError)

	nodeAddr := fmt.Sprintf("http://127.0.0.1:%d", spec.Node.WebInterfacePort)
	addNodeAddrFlag := func(cmd *flag.FlagSet) {
		cmd.StringVar(&nodeAddr, "node", nodeAddr, "HTTP API `ADDRESS` of cxchain node")
		cmd.StringVar(&nodeAddr, "n", nodeAddr, "shorthand for 'node'")
	}

	cmdPrelude := func(args []string, argsName string) *flag.FlagSet {
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

	cmdMap := cxutil.NewCommandMap(rootCmd, 3, cxutil.DefaultUsageSignature()).
		Add("connect", func(args []string) {
			cmd := cmdPrelude(args, "addresses")
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
		}).
		Add("disconnect", func(args []string) {
			cmd := cmdPrelude(args, "conn_ids")
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
		}).
		Add("list", func(args []string) {
			fmt.Println("TODO @evanlinjin: Not implemented.")
		})

	os.Exit(cmdMap.ParseAndRun(args[1:]))
}
