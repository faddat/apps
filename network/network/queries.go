package network

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ignite/cli/ignite/pkg/cosmoserror"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/pkg/errors"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	projecttypes "github.com/tendermint/spn/x/project/types"
	rewardtypes "github.com/tendermint/spn/x/reward/types"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/apps/network/network/networktypes"
)

// ErrObjectNotFound is returned when the query returns a not found error.
var ErrObjectNotFound = errors.New("query object not found")

// ChainLaunch fetches the chain launch from Network by launch id.
func (n Network) ChainLaunch(ctx context.Context, id uint64) (networktypes.ChainLaunch, error) {
	n.ev.Send("Fetching chain information", events.ProgressStart())

	res, err := n.launchQuery.
		Chain(ctx,
			&launchtypes.QueryGetChainRequest{
				LaunchID: id,
			},
		)
	if err != nil {
		return networktypes.ChainLaunch{}, err
	}

	return networktypes.ToChainLaunch(res.Chain), nil
}

// ChainLaunchesWithReward fetches the chain launches with rewards from Network.
func (n Network) ChainLaunchesWithReward(ctx context.Context, pagination *query.PageRequest) ([]networktypes.ChainLaunch, error) {
	g, ctx := errgroup.WithContext(ctx)

	n.ev.Send("Fetching chains information", events.ProgressStart())
	res, err := n.launchQuery.
		ChainAll(ctx, &launchtypes.QueryAllChainRequest{
			Pagination: pagination,
		})
	if err != nil {
		return nil, err
	}

	n.ev.Send("Fetching reward information", events.ProgressUpdate())
	var chainLaunches []networktypes.ChainLaunch
	var mu sync.Mutex

	// Parse fetched chains and fetch rewards
	for _, chain := range res.Chain {
		chain := chain
		g.Go(func() error {
			chainLaunch := networktypes.ToChainLaunch(chain)
			reward, err := n.ChainReward(ctx, chain.LaunchID)
			if err != nil && !errors.Is(err, ErrObjectNotFound) {
				return err
			}
			chainLaunch.Reward = reward.RemainingCoins.String()
			mu.Lock()
			chainLaunches = append(chainLaunches, chainLaunch)
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	// sort filenames by launch id
	sort.Slice(chainLaunches, func(i, j int) bool {
		return chainLaunches[i].ID > chainLaunches[j].ID
	})
	return chainLaunches, nil
}

// GenesisInformation returns all the information to construct the genesis from a chain ID.
func (n Network) GenesisInformation(ctx context.Context, launchID uint64) (gi networktypes.GenesisInformation, err error) {
	genAccs, err := n.GenesisAccounts(ctx, launchID)
	if err != nil {
		return gi, fmt.Errorf("error querying genesis accounts: %w", err)
	}

	vestingAccs, err := n.VestingAccounts(ctx, launchID)
	if err != nil {
		return gi, fmt.Errorf("error querying vesting accounts: %w", err)
	}

	genVals, err := n.GenesisValidators(ctx, launchID)
	if err != nil {
		return gi, fmt.Errorf("error querying genesis validators: %w", err)
	}

	paramChanges, err := n.ParamChanges(ctx, launchID)
	if err != nil {
		return gi, fmt.Errorf("error querying param changes: %w", err)
	}

	return networktypes.NewGenesisInformation(genAccs, vestingAccs, genVals, paramChanges), nil
}

// GenesisAccounts returns the list of approved genesis accounts for a launch from SPN.
func (n Network) GenesisAccounts(ctx context.Context, launchID uint64) (genAccs []networktypes.GenesisAccount, err error) {
	n.ev.Send("Fetching genesis accounts", events.ProgressStart())
	res, err := n.launchQuery.
		GenesisAccountAll(ctx,
			&launchtypes.QueryAllGenesisAccountRequest{
				LaunchID: launchID,
			},
		)
	if err != nil {
		return genAccs, err
	}

	for _, acc := range res.GenesisAccount {
		genAccs = append(genAccs, networktypes.ToGenesisAccount(acc))
	}

	return genAccs, nil
}

// VestingAccounts returns the list of approved genesis vesting accounts for a launch from SPN.
func (n Network) VestingAccounts(ctx context.Context, launchID uint64) (vestingAccs []networktypes.VestingAccount, err error) {
	n.ev.Send("Fetching genesis vesting accounts", events.ProgressStart())
	res, err := n.launchQuery.
		VestingAccountAll(ctx,
			&launchtypes.QueryAllVestingAccountRequest{
				LaunchID: launchID,
			},
		)
	if err != nil {
		return vestingAccs, err
	}

	for i, acc := range res.VestingAccount {
		parsedAcc, err := networktypes.ToVestingAccount(acc)
		if err != nil {
			return vestingAccs, errors.Wrapf(err, "error parsing vesting account %d", i)
		}

		vestingAccs = append(vestingAccs, parsedAcc)
	}

	return vestingAccs, nil
}

// GenesisValidators returns the list of approved genesis validators for a launch from SPN.
func (n Network) GenesisValidators(ctx context.Context, launchID uint64) (genVals []networktypes.GenesisValidator, err error) {
	n.ev.Send("Fetching genesis validators", events.ProgressStart())
	res, err := n.launchQuery.
		GenesisValidatorAll(ctx,
			&launchtypes.QueryAllGenesisValidatorRequest{
				LaunchID: launchID,
			},
		)
	if err != nil {
		return genVals, err
	}

	for _, acc := range res.GenesisValidator {
		genVals = append(genVals, networktypes.ToGenesisValidator(acc))
	}

	return genVals, nil
}

// ParamChanges returns the list of approved param changes for a launch from SPN.
func (n Network) ParamChanges(ctx context.Context, launchID uint64) (paramChanges []networktypes.ParamChange, err error) {
	n.ev.Send("Fetching param changes", events.ProgressStart())
	res, err := n.launchQuery.
		ParamChangeAll(ctx,
			&launchtypes.QueryAllParamChangeRequest{
				LaunchID: launchID,
			},
		)
	if err != nil {
		return paramChanges, err
	}

	for _, pc := range res.ParamChanges {
		paramChanges = append(paramChanges, networktypes.ToParamChange(pc))
	}

	return paramChanges, nil
}

// MainnetAccount returns the project mainnet account for a launch from SPN.
func (n Network) MainnetAccount(
	ctx context.Context,
	projectID uint64,
	address string,
) (acc networktypes.MainnetAccount, err error) {
	n.ev.Send(
		fmt.Sprintf("Fetching project %d mainnet account %s", projectID, address),
		events.ProgressStart(),
	)
	res, err := n.projectQuery.
		MainnetAccount(ctx,
			&projecttypes.QueryGetMainnetAccountRequest{
				ProjectID: projectID,
				Address:   address,
			},
		)
	if errors.Is(cosmoserror.Unwrap(err), cosmoserror.ErrNotFound) {
		return networktypes.MainnetAccount{}, ErrObjectNotFound
	} else if err != nil {
		return acc, err
	}

	return networktypes.ToMainnetAccount(res.MainnetAccount), nil
}

// MainnetAccounts returns the list of project mainnet accounts for a launch from SPN.
func (n Network) MainnetAccounts(ctx context.Context, projectID uint64) (genAccs []networktypes.MainnetAccount, err error) {
	n.ev.Send("Fetching project mainnet accounts", events.ProgressStart())
	res, err := n.projectQuery.
		MainnetAccountAll(ctx,
			&projecttypes.QueryAllMainnetAccountRequest{
				ProjectID: projectID,
			},
		)
	if err != nil {
		return genAccs, err
	}

	for _, acc := range res.MainnetAccount {
		genAccs = append(genAccs, networktypes.ToMainnetAccount(acc))
	}

	return genAccs, nil
}

func (n Network) GenesisAccount(ctx context.Context, launchID uint64, address string) (networktypes.GenesisAccount, error) {
	n.ev.Send("Fetching genesis accounts", events.ProgressStart())
	res, err := n.launchQuery.GenesisAccount(ctx, &launchtypes.QueryGetGenesisAccountRequest{
		LaunchID: launchID,
		Address:  address,
	})
	if errors.Is(cosmoserror.Unwrap(err), cosmoserror.ErrNotFound) {
		return networktypes.GenesisAccount{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.GenesisAccount{}, err
	}
	return networktypes.ToGenesisAccount(res.GenesisAccount), nil
}

func (n Network) VestingAccount(ctx context.Context, launchID uint64, address string) (networktypes.VestingAccount, error) {
	n.ev.Send("Fetching vesting accounts", events.ProgressStart())
	res, err := n.launchQuery.VestingAccount(ctx, &launchtypes.QueryGetVestingAccountRequest{
		LaunchID: launchID,
		Address:  address,
	})
	if errors.Is(cosmoserror.Unwrap(err), cosmoserror.ErrNotFound) {
		return networktypes.VestingAccount{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.VestingAccount{}, err
	}
	return networktypes.ToVestingAccount(res.VestingAccount)
}

func (n Network) GenesisValidator(ctx context.Context, launchID uint64, address string) (networktypes.GenesisValidator, error) {
	n.ev.Send("Fetching genesis validator", events.ProgressStart())
	res, err := n.launchQuery.GenesisValidator(ctx, &launchtypes.QueryGetGenesisValidatorRequest{
		LaunchID: launchID,
		Address:  address,
	})
	if errors.Is(cosmoserror.Unwrap(err), cosmoserror.ErrNotFound) {
		return networktypes.GenesisValidator{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.GenesisValidator{}, err
	}
	return networktypes.ToGenesisValidator(res.GenesisValidator), nil
}

// ChainReward fetches the chain reward from SPN by launch id.
func (n Network) ChainReward(ctx context.Context, launchID uint64) (rewardtypes.RewardPool, error) {
	res, err := n.rewardQuery.
		RewardPool(ctx,
			&rewardtypes.QueryGetRewardPoolRequest{
				LaunchID: launchID,
			},
		)

	if errors.Is(cosmoserror.Unwrap(err), cosmoserror.ErrNotFound) {
		return rewardtypes.RewardPool{}, ErrObjectNotFound
	} else if err != nil {
		return rewardtypes.RewardPool{}, err
	}
	return res.RewardPool, nil
}

// ChainID fetches the network chain id.
func (n Network) ChainID(ctx context.Context) (string, error) {
	status, err := n.cosmos.Status(ctx)
	if err != nil {
		return "", err
	}
	return status.NodeInfo.Network, nil
}
