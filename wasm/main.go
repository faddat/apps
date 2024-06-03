package main

import (
	"context"
	"os"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/wasm/cmd"
)

type App struct{}

func (App) Manifest(context.Context) (*plugin.Manifest, error) {
	m := &plugin.Manifest{Name: "wasm"}
	m.ImportCobraCommand(cmd.NewWasm(), "ignite")
	return m, nil
}

func (App) Execute(_ context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// Run the "hermes" command as if it were a root command. To do
	// so remove the first two arguments which are "ignite relayer"
	// from OSArgs to treat "hermes" as the root command.
	os.Args = c.OsArgs[1:]
	return cmd.NewWasm().Execute()
}

func (App) ExecuteHookPre(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (App) ExecuteHookPost(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func (App) ExecuteHookCleanUp(context.Context, *plugin.ExecutedHook, plugin.ClientAPI) error {
	return nil
}

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins: map[string]hplugin.Plugin{
			"wasm": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
