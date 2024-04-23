package cmd

import (
	"context"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/ignite/gex/pkg/xurl"
	"github.com/ignite/gex/services/explorer"
)

// ExecuteGex executes explorer gex subcommand.
func ExecuteGex(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	rpcAddress, err := flags.GetString(flagRPCAddress)
	if err != nil {
		return errors.Errorf("could not get --%s flag: %s", flagRPCAddress, err)
	}

	hostURL, err := xurl.Parse(rpcAddress)
	if err != nil {
		return err
	}

	return explorer.Run(ctx, hostURL.String())
}
