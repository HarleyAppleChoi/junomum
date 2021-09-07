package utils

import (
	"fmt"
	"strings"

	"github.com/forbole/flowJuno/client"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"

	"github.com/rs/zerolog/log"

	db "github.com/forbole/flowJuno/db/postgresql"

	"github.com/forbole/flowJuno/types"
)

/*
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
} */

// --------------------------------------------------------------------------------------------------------------------

// GetAccounts returns the account data for the given addresses
func GetAccounts(addresses []string, height int64, client client.Proxy) ([]types.Account, error) {
	log.Debug().Str("module", "auth").Str("operation", "accounts").Str("Height",string(rune(height))).Msg("getting accounts data")
	var accounts []types.Account

	for _, address := range addresses {
		fmt.Println("GetAccounts:"+address)
		if address==""{
			continue
		}
		//not working atm because of flow bug
		//account,err:=client.Client().GetAccountAtBlockHeight(client.Ctx(),flow.HexToAddress(address),uint64(height))

		account,err:=client.Client().GetAccount(client.Ctx(),flow.HexToAddress(address))

		if err != nil {
			return nil,err
		}

		if account == nil {
			return nil, fmt.Errorf("address is not valid and cannot get details")
		}
		
		accounts = append(accounts, types.NewAccount(account.Address.String()))

	}

	return accounts, nil
}

// UpdateAccounts takes the given addresses and for each one queries the chain
// retrieving the account data and stores it inside the database.
func UpdateAccounts(addresses []string, db *db.Db, height int64, client client.Proxy) error {
	accounts, err := GetAccounts(addresses, height, client)
	if err != nil {
		return err
	}

	lockedAccount,err:=GetLockedTokenAccount(addresses, height, client)
	if err!=nil{
		return err
	}

	err = db.SaveAccounts(accounts)
	if err!=nil{
		return err
	}

	return db.SaveLockedTokenAccounts(lockedAccount)


}

func GetLockedTokenAccount(addresses []string, height int64, client client.Proxy)([]types.LockedAccount,error){
	script:=fmt.Sprintf(`
	import LockedTokens from %s

	pub fun main(account: Address): Address {
	
		let lockedAccountInfoRef = getAccount(account)
			.getCapability<&LockedTokens.TokenHolder{LockedTokens.LockedAccountInfo}>(
				LockedTokens.LockedAccountInfoPublicPath
			)
			.borrow()
			?? panic("Could not borrow a reference to public LockedAccountInfo")
	
		return lockedAccountInfoRef.getLockedAccountAddress()
	}

	`,client.Contract().LockedTokens)

	var lockedAccount []types.LockedAccount

	for _,address:=range addresses{
		if address==""{
			continue
		}
		flowAddress:=flow.HexToAddress(address)
		candanceAddress:=cadence.Address(flowAddress)
		//val,err:=cadence.NewValue(candanceAddress)
		candenceArr:=[]cadence.Value{candanceAddress}

		catchError:=`Could not borrow a reference to public LockedAccountInfo`
		value,err:=client.Client().ExecuteScriptAtLatestBlock(client.Ctx(),[]byte(script),candenceArr)
		if err==nil{
			fmt.Println("LockedAccountGet!"+value.String())

			lockedAccount=append(lockedAccount,types.NewLockedAccount(address,value.String()))	
			
		}else if (strings.Contains(err.Error(),catchError)){
			//This account don't have a locked account
			continue
		}else if err!=nil{
			return nil,err
		}
	}
	return lockedAccount,nil
}