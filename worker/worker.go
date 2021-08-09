package worker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/desmos-labs/juno/types/logging"
	"github.com/onflow/flow-go-sdk"

	tmjson "github.com/tendermint/tendermint/libs/json"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/desmos-labs/juno/modules"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/types"
)

// Worker defines a job consumer that is responsible for getting and
// aggregating block and associated data and exporting it to a database.
type Worker struct {
	queue          types.HeightQueue
	encodingConfig *params.EncodingConfig
	cp             *client.Proxy
	db             db.Database
	modules        []modules.Module
}

// NewWorker allows to create a new Worker implementation.
func NewWorker(config *Config) Worker {
	return Worker{
		encodingConfig: config.EncodingConfig,
		cp:             config.ClientProxy,
		queue:          config.Queue,
		db:             config.Database,
		modules:        config.Modules,
	}
}

// Start starts a worker by listening for new jobs (block heights) from the
// given worker queue. Any failed job is logged and re-enqueued.
func (w Worker) Start() {
	for i := range w.queue {
		if err := w.process(i); err != nil {
			// re-enqueue any failed job
			// TODO: Implement exponential backoff or max retries for a block height.
			go func() {
				log.Error().Err(err).Int64("height", i).Msg("re-enqueueing failed block")
				w.queue <- i
			}()
		}
	}
}

// process defines the job consumer workflow. It will fetch a block for a given
// height and associated metadata and export it to a database. It returns an
// error if any export process fails.
func (w Worker) process(height int64) error {
	exists, err := w.db.HasBlock(height)
	if err != nil {
		return err
	}

	if exists {
		log.Debug().Int64("height", height).Msg("skipping already exported block")
		return nil
	}
/* 
	if height == 0 {
		cfg := types.Cfg.GetParsingConfig()
		var genesis *tmtypes.GenesisDoc
		if strings.TrimSpace(cfg.GetGenesisFilePath()) != "" {
			genesis, err = w.getGenesisFromFilePath(cfg.GetGenesisFilePath())
			if err != nil {
				return err
			}
		} else {
			genesis, err = w.getGenesisFromRPC()
			if err != nil {
				return err
			}
		}

		return w.HandleGenesis(genesis)
	}
 */
	//log.Debug().Int64("height", height).Msg("processing block")

	block, err := w.cp.Block(height)
	if err != nil {
		log.Error().Err(err).Int64("height", height).Msg("failed to get block")
		return err
	}

	txs, err := w.cp.Txs(block)
	if err != nil {
		log.Error().Err(err).Int64("height", height).Msg("failed to get transaction Result for block")
		return err
	}

	vals, err := w.cp.NodeOperators(height)
	if err != nil {
		log.Error().Err(err).Int64("height", height).Msg("failed to get node operators for block")
		return err
	}

	return w.ExportBlock(block, txs, vals)
}

// getGenesisFromRPC returns the genesis read from the RPC endpoint
func (w Worker) getGenesisFromRPC() (*tmtypes.GenesisDoc, error) {
	log.Debug().Msg("getting genesis")
	response, err := w.cp.Genesis()
	if err != nil {
		log.Error().Err(err).Msg("failed to get genesis")
		return nil, err
	}
	return response.Genesis, nil
}

// getGenesisFromFilePath tries reading the genesis doc from the given path
func (w Worker) getGenesisFromFilePath(path string) (*tmtypes.GenesisDoc, error) {
	log.Debug().Str("path", path).Msg("reading genesis from file")

	bz, err := tmos.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Msg("failed to read genesis file")
		return nil, err
	}

	var genDoc tmtypes.GenesisDoc
	err = tmjson.Unmarshal(bz, &genDoc)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal genesis doc")
		return nil, err
	}

	return &genDoc, nil
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
// in the order in which they have been registered.
func (w Worker) HandleGenesis(genesis *tmtypes.GenesisDoc) error {
	var appState map[string]json.RawMessage
	if err := json.Unmarshal(genesis.AppState, &appState); err != nil {
		return fmt.Errorf("error unmarshalling genesis doc %s: %s", appState, err.Error())
	}

	// Call the genesis handlers
	for _, module := range w.modules {
		if genesisModule, ok := module.(modules.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(genesis, appState); err != nil {
				logging.LogGenesisError(module, err)
			}
		}
	}

	return nil
}

// SaveValidators persists a list of Tendermint validators with an address and a
// consensus public key. An error is returned if the public key cannot be Bech32
// encoded or if the DB write fails.
func (w Worker) SaveValidators(vals []*tmtypes.Validator) error {
	err := w.db.SaveValidators(validators)
	if err != nil {
		return fmt.Errorf("error while saving validators: %s", err)
	}

	return nil
}

func (w Worker) SaveNodeInfos(vals []*types.NodeInfo)error{
	err:=w.db.SaveNodeInfos(vals)
	if err!=nil{
		return fmt.Errorf("error while saving node infos: %s", err)
	}
	return nil
}

// ExportBlock accepts a finalized block and a corresponding set of transactions
// and persists them to the database along with attributable metadata. An error
// is returned if the write fails.
func (w Worker) ExportBlock(b *flow.Block, txs []*types.Txs, vals *types.NodeOperators) error {
	// Save all validators
	err := w.SaveNodeInfos(vals.NodeInfos)
	if err != nil {
		return err
	}

	// Save the block
	err = w.db.SaveBlock(b)
	if err != nil {
		log.Error().Err(err).Int64("height", int64(b.BlockHeader.Height)).Msg("failed to persist block")
		return err
	}


	// Call the block handlers
	for _, module := range w.modules {
		if blockModule, ok := module.(modules.BlockModule); ok {
			err = blockModule.HandleBlock(b, txs, vals)
			if err != nil {
				logging.LogBLockError(module, b, err)
			}
		}
	}

	// Export the transactions
	return w.ExportTxEvents(txs)
}

// ExportCommit accepts a block commitment and a corresponding set of
// validators for the commitment and persists them to the database. An error is
// returned if any write fails or if there is any missing aggregated data.
func (w Worker) ExportCommit(commit *tmtypes.Commit, vals *tmctypes.ResultValidators) error {
	var signatures []*types.CommitSig
	for _, commitSig := range commit.Signatures {
		// Avoid empty commits
		if commitSig.Signature == nil {
			continue
		}

		valAddr := sdk.ConsAddress(commitSig.ValidatorAddress)
		val := findValidatorByAddr(valAddr.String(), vals)
		if val == nil {
			err := fmt.Errorf("failed to find validator")
			log.Error().
				Err(err).
				Int64("height", commit.Height).
				Str("validator_hex", commitSig.ValidatorAddress.String()).
				Str("validator_bech32", valAddr.String()).
				Time("commit_timestamp", commitSig.Timestamp).
				Send()
			return err
		}

		signatures = append(signatures, types.NewCommitSig(
			types.ConvertValidatorAddressToBech32String(commitSig.ValidatorAddress),
			val.VotingPower,
			val.ProposerPriority,
			commit.Height,
			commitSig.Timestamp,
		))
	}

	err := w.db.SaveCommitSignatures(signatures)
	if err != nil {
		return fmt.Errorf("error while saving commit signatures: %s", err)
	}

	return nil
}

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if the write fails.
func (w Worker) ExportTxEvents(txs []*types.Txs) error {
	// Handle all the transactions inside the block
	for _, tx := range txs {
		// Save the transaction itself
		err := w.db.SaveTx(tx)
		if err != nil {
			log.Error().Err(err).Str("Height", string(tx.Height)).Msg("failed to handle transaction")
			return err
		}

		for _,event := range tx.Events{
			transaction,err:=w.cp.GetTransaction(event.TransactionID.String())
			if err!=nil{
				return err
			}
			transaction.

						
		}

		// Call the tx handlers
		for _, module := range w.modules {
			if transactionModule, ok := module.(modules.TransactionModule); ok {
				err = transactionModule.HandleTx(tx)
				if err != nil {
					logging.LogTxError(module, tx, err)
				}
			}
		}

		// Handle all the messages contained inside the transaction
		for i, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg
			err = w.encodingConfig.Marshaler.UnpackAny(msg, &stdMsg)
			if err != nil {
				return err
			}

			// Call the handlers
			for _, module := range w.modules {
				if messageModule, ok := module.(modules.MessageModule); ok {
					err = messageModule.HandleMsg(i, stdMsg, tx)
					if err != nil {
						logging.LogMsgError(module, tx, stdMsg, err)
					}
				}
			}
		}
	}

	return nil
}
