package main

import (
	"context"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"chain-info/cmd"
)

type App struct{}

func (App) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "chain-info",
		Commands: cmd.GetCommands(),
	}, nil
}

func (App) Execute(ctx context.Context, c *plugin.ExecutedCommand, api plugin.ClientAPI) error {
	// Remove the first two elements "ignite" and "flags" from OsArgs.
	args := c.OsArgs[2:]

	chainInfo, err := api.GetChainInfo(ctx)
	if err != nil {
		return errors.Errorf("failed to get chain info: %s", err)
	}

	ch, err := chain.New(chainInfo.AppPath)
	if err != nil {
		return errors.Errorf("failed to create a new chain object from app path: %s", err)
	}

	switch args[0] {
	case "info":
		return cmd.ExecuteInfo(ctx, c, ch)
	case "build":
		return cmd.ExecuteBuild(ctx, c, ch)
	default:
		return errors.Errorf("unknown command: %s", c.Path)
	}
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
			"chain-info": plugin.NewGRPC(&App{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
