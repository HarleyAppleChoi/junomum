package messages

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/HarleyAppleChoi/junomum/db"
	"github.com/HarleyAppleChoi/junomum/modules/modules"
	"github.com/HarleyAppleChoi/junomum/types"
)

var _ modules.Module = &Module{}

// Module represents the module allowing to store messages properly inside a dedicated table
type Module struct {
	parser MessageAddressesParser
	cdc    codec.Marshaler
	db     db.Database
}

func NewModule(parser MessageAddressesParser, cdc codec.Marshaler, db db.Database) *Module {
	return &Module{
		parser: parser,
		cdc:    cdc,
		db:     db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "messages"
}

// HandleEvent implements modules.MessageModule
func (m *Module) HandleEvent(index int, msg sdk.Msg, tx *types.Txs) error {
	//return HandleEvent(index, msg, tx, m.parser, m.cdc, m.db)
	return nil
}
