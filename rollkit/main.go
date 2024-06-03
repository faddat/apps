package main

import (
	"context"
	"os"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/rollkit/cmd"
)

var _ plugin.Interface = app{}

type App struct{}

func (App) Manifest(context.Context) (*plugin.Manifest, error) {
	m := &plugin.Manifest{Name: "rollkit"}
	m.ImportCobraCommand(cmd.NewRollkit(), "ignite")
	return m, nil
}

func (App) Execute(_ context.Context, c *plugin.ExecutedCommand, _ plugin.ClientAPI) error {
	// Instead of a switch on c.Use, we run the root command like if
	// we were in a command line context. This implies to set os.Args
	// correctly.
	// Remove the first arg "ignite" from OSArgs because our rollkit
	// command root is "rollkit" not "ignite".
	os.Args = c.OsArgs[1:]
	return cmd.NewRollkit().Execute()
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
			"rollkit": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
