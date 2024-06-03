package main

import (
	"fmt"
	"os"

	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/spf13/cobra"

	chaininfo "github.com/ignite/apps/examples/chain-info"
	chaininfocmd "github.com/ignite/apps/examples/chain-info/cmd"
	flags "github.com/ignite/apps/examples/flags"
	flagscmd "github.com/ignite/apps/examples/flags/cmd"
	healthmonitor "github.com/ignite/apps/examples/health-monitor"
	healthmonitorcmd "github.com/ignite/apps/examples/health-monitor/cmd"
	helloworld "github.com/ignite/apps/examples/hello-world"
	helloworldcmd "github.com/ignite/apps/examples/hello-world/cmd"
	hooks "github.com/ignite/apps/examples/hooks"
	hookscmd "github.com/ignite/apps/examples/hooks/cmd"
	explorer "github.com/ignite/apps/explorer"
	explorercmd "github.com/ignite/apps/explorer/cmd"
	hermescmd "github.com/ignite/apps/hermes/cmd"
	marketplacecmd "github.com/ignite/apps/marketplace/cmd"
	wasmcmd "github.com/ignite/apps/wasm/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "apps",
		Short: "debug apps commands",
	}

	// Add apps with ignite app commands.
	newCmdFromApp(rootCmd, chaininfo.App{}, chaininfocmd.GetCommands())
	newCmdFromApp(rootCmd, flags.App{}, flagscmd.GetCommands())
	newCmdFromApp(rootCmd, healthmonitor.App{}, healthmonitorcmd.GetCommands())
	newCmdFromApp(rootCmd, helloworld.App{}, helloworldcmd.GetCommands())
	newCmdFromApp(rootCmd, hooks.App{}, hookscmd.GetCommands())
	newCmdFromApp(rootCmd, explorer.App{}, explorercmd.GetCommands())
	// Add ignite app commands for debugging here.

	// Add apps with cobra commands.
	rootCmd.AddCommand(
		hermescmd.NewRelayer(),
		marketplacecmd.NewMarketplace(),
		wasmcmd.NewWasm(),
		// Add cobra commands for debugging here.
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCmdFromApp(rootCmd *cobra.Command, app plugin.Interface, commands []*plugin.Command) {
	for _, cmd := range commands {
		cobraCmd, err := cmd.ToCobraCommand()
		if err != nil {
			panic(err)
		}

		newCmdFromApp(cobraCmd, app, cmd.GetCommands())
		rootCmd.AddCommand(cobraCmd)
	}
}
