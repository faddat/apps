package main

import (
	"context"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"hello-world/cmd"
)

type App struct{}

func (App) Manifest(context.Context) (*plugin.Manifest, error) {
	return &plugin.Manifest{
		Name:     "hello-world",
		Commands: cmd.GetCommands(),
	}, nil
}

func (App) Execute(context.Context, *plugin.ExecutedCommand, plugin.ClientAPI) error {
	fmt.Println("Hello, world!")
	return nil
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
			"hello-world": plugin.NewGRPC(&app{}),
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
