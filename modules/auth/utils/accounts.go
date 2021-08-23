package utils

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/forbole/flowJuno/client"

	"github.com/rs/zerolog/log"

	"github.com/forbole/flowJuno/db"
	"github.com/forbole/flowJuno/types"

	"github.com/forbole/flowJuno/types"
)

// GetGenesisAccounts parses the given appState and returns the genesis accounts
func GetGenesisAccounts(appState map[string]json.RawMessage, cdc codec.Marshaler) ([]types.Account, error) {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	var authState authtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[authtypes.ModuleName], &authState); err != nil {
		return nil, err
	}

	// Store the accounts
	accounts := make([]types.Account, len(authState.Accounts))
	for index, account := range authState.Accounts {
		var accountI authtypes.AccountI
		err := cdc.UnpackAny(account, &accountI)
		if err != nil {
			return nil, err
		}

		accounts[index] = types.NewAccount(accountI.GetAddress().String(), accountI)
	}

	return accounts, nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetAccounts returns the account data for the given addresses
func GetAccounts(addresses []string, cdc codec.Marshaler, height int64, client client.Proxy) ([]types.Account, error) {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Msg("getting accounts data")
	var accounts []types.Account
	header := client.GetHeightRequestHeader(height)

	for _, address := range addresses {
		client.GetAccount(ctx,)
		if err != nil {
			log.Error().Str("module", "auth").Err(err).Int64("height", height).
				Str("Auth", "Get Account").Msg("error while getting accounts")
		}
		if res == nil {
			return nil, fmt.Errorf("address is not valid and cannot get details")
		}
		var accountI authtypes.AccountI
		err = cdc.UnpackAny(res.Account, &accountI)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, types.NewAccount(accountI.GetAddress().String(), accountI))

	}

	return accounts, nil
}

// UpdateAccounts takes the given addresses and for each one queries the chain
// retrieving the account data and stores it inside the database.
func UpdateAccounts(addresses []string, cdc codec.Marshaler, db *db.Database, height int64, client client.Client) error {
	accounts, err := GetAccounts(addresses, cdc, height, client)
	if err != nil {
		return err
	}

	return db.SaveAccounts(accounts)
}
