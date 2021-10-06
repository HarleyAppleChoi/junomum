package staking

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/forbole/flowJuno/types"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"

	"github.com/forbole/flowJuno/client"
	"github.com/forbole/flowJuno/modules/utils"

	database "github.com/forbole/flowJuno/db/postgresql"
)

func getDelegatorCommitted(nodeId string, delegatorID uint32, block *flow.Block, db *database.Db, flowClient client.Proxy) error {
	log.Trace().Str("module", "staking").Int64("height", int64(block.Height)).
		Msg("updating node unstaking tokens")
	script := fmt.Sprintf(`
	import FlowIDTableStaking from %s
	pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
	  let delInfo = FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
	  return delInfo.tokensCommitted
  }`, flowClient.Contract().StakingTable)

	args := []cadence.Value{cadence.NewString(nodeId), cadence.NewUInt32(delegatorID)}
	value, err := flowClient.Client().ExecuteScriptAtLatestBlock(flowClient.Ctx(), []byte(script), args)
	if err != nil {
		return err
	}

	committed, err := utils.CadenceConvertUint64(value)
	if err != nil {
		return err
	}

	return db.SaveDelegatorCommitted(types.NewDelegatorCommitted(committed, block.Height, nodeId, delegatorID))
}

func getDelegatorInfo(nodeId string, delegatorID uint32, block *flow.Block, db *database.Db, flowClient client.Proxy) error {
	log.Trace().Str("module", "staking").Int64("height", int64(block.Height)).
		Msg("updating node unstaking tokens")
	script := fmt.Sprintf(`
	import FlowIDTableStaking from %s
	pub fun main(nodeID: String, delegatorID: UInt32): FlowIDTableStaking.DelegatorInfo {
	  return FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
  }`, flowClient.Contract().StakingTable)

	args := []cadence.Value{cadence.NewString(nodeId), cadence.NewUInt32(delegatorID)}
	value, err := flowClient.Client().ExecuteScriptAtLatestBlock(flowClient.Ctx(), []byte(script), args)
	if err != nil {
		return err
	}

	committed, err := types.DelegatorNodeInfoFromCadence(value)
	if err != nil {
		return err
	}

	return db.SaveDelegatorInfo(types.NewDelegatorInfo(committed, block.Height, nodeId, delegatorID))
}

func getDelegatorRequest(nodeId string, delegatorID uint32, block *flow.Block, db *database.Db, flowClient client.Proxy) error {
	log.Trace().Str("module", "staking").Int64("height", int64(block.Height)).
		Msg("updating node unstaking tokens")
	script := fmt.Sprintf(`
	import FlowIDTableStaking from %s
	pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
	  let delInfo = FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
	  return delInfo.tokensRequestedToUnstake
  }`, flowClient.Contract().StakingTable)

	args := []cadence.Value{cadence.NewString(nodeId), cadence.NewUInt32(delegatorID)}
	value, err := flowClient.Client().ExecuteScriptAtLatestBlock(flowClient.Ctx(), []byte(script), args)
	if err != nil {
		return err
	}

	committed, err := utils.CadenceConvertUint64(value)
	if err != nil {
		return err
	}

	return db.SaveDelegatorRequest(types.NewDelegatorRequest(committed, int64(block.Height), nodeId, delegatorID))
}

func DelegatorRewarded(nodeId string, delegatorID uint32, block *flow.Block, db *database.Db, flowClient client.Proxy) error {
	log.Trace().Str("module", "staking").Int64("height", int64(block.Height)).
		Msg("updating node unstaking tokens")
	script := fmt.Sprintf(`
	import FlowIDTableStaking from %s
	pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
	  let delInfo = FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
	  return delInfo.tokensRewarded
  }`, flowClient.Contract().StakingTable)

	args := []cadence.Value{cadence.NewString(nodeId), cadence.NewUInt32(delegatorID)}
	value, err := flowClient.Client().ExecuteScriptAtLatestBlock(flowClient.Ctx(), []byte(script), args)
	if err != nil {
		return err
	}

	committed, err := utils.CadenceConvertUint64(value)
	if err != nil {
		return err
	}

	return db.SaveDelegatorRewarded(types.NewDelegatorRewarded(committed, int64(block.Height), nodeId, delegatorID))
}

func getDelegatorStaked(nodeId string, delegatorID uint32, block *flow.Block, db *database.Db, flowClient client.Proxy) error {
	log.Trace().Str("module", "staking").Int64("height", int64(block.Height)).
		Msg("updating node unstaking tokens")
	script := fmt.Sprintf(`
	import FlowIDTableStaking from %s
	pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
	  let delInfo = FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
	  return delInfo.tokensStaked
  }`, flowClient.Contract().StakingTable)

	args := []cadence.Value{cadence.NewString(nodeId), cadence.NewUInt32(delegatorID)}
	value, err := flowClient.Client().ExecuteScriptAtLatestBlock(flowClient.Ctx(), []byte(script), args)
	if err != nil {
		return err
	}

	staked, err := utils.CadenceConvertUint64(value)
	if err != nil {
		return err
	}

	return db.SaveDelegatorStaked(types.NewDelegatorStaked(staked, int64(block.Height), nodeId, delegatorID))
}

func getDelegatorUnstaked(nodeId string, delegatorID uint32, block *flow.Block, db *database.Db, flowClient client.Proxy) error {
	log.Trace().Str("module", "staking").Int64("height", int64(block.Height)).
		Msg("updating node unstaking tokens")
	script := fmt.Sprintf(`
	import FlowIDTableStaking from %s
	pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
	  let delInfo = FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
	  return delInfo.tokensUnstaked
  }`, flowClient.Contract().StakingTable)

	args := []cadence.Value{cadence.NewString(nodeId), cadence.NewUInt32(delegatorID)}
	value, err := flowClient.Client().ExecuteScriptAtLatestBlock(flowClient.Ctx(), []byte(script), args)
	if err != nil {
		return err
	}

	staked, err := utils.CadenceConvertUint64(value)
	if err != nil {
		return err
	}

	return db.SaveDelegatorUnstaked(types.NewDelegatorUnstaked(staked, int64(block.Height), nodeId, delegatorID))
}

func getDelegatorUnstaking(nodeId string, delegatorID uint32, block *flow.Block, db *database.Db, flowClient client.Proxy) error {
	log.Trace().Str("module", "staking").Int64("height", int64(block.Height)).
		Msg("updating node unstaking tokens")
	script := fmt.Sprintf(`
	import FlowIDTableStaking from %s
	pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
	  let delInfo = FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
	  return delInfo.tokensUnstaking
  }`, flowClient.Contract().StakingTable)

	args := []cadence.Value{cadence.NewString(nodeId), cadence.NewUInt32(delegatorID)}
	value, err := flowClient.Client().ExecuteScriptAtLatestBlock(flowClient.Ctx(), []byte(script), args)
	if err != nil {
		return err
	}

	staking, err := utils.CadenceConvertUint64(value)
	if err != nil {
		return err
	}

	return db.SaveDelegatorUnstaking(types.NewDelegatorUnstaking(staking, int64(block.Height), nodeId, delegatorID))
}

func getDelegatorUnstakingRequest(nodeId string, delegatorID uint32, block *flow.Block, db *database.Db, flowClient client.Proxy) error {
	log.Trace().Str("module", "staking").Int64("height", int64(block.Height)).
		Msg("updating node unstaking tokens")
	script := fmt.Sprintf(`
	import FlowIDTableStaking from %s
	pub fun main(nodeID: String, delegatorID: UInt32): UFix64 {
	  let delInfo = FlowIDTableStaking.DelegatorInfo(nodeID: nodeID, delegatorID: delegatorID)
	  return delInfo.tokensRequestedToUnstake
  }`, flowClient.Contract().StakingTable)

	args := []cadence.Value{cadence.NewString(nodeId), cadence.NewUInt32(delegatorID)}
	value, err := flowClient.Client().ExecuteScriptAtLatestBlock(flowClient.Ctx(), []byte(script), args)
	if err != nil {
		return err
	}

	staking, err := utils.CadenceConvertUint64(value)
	if err != nil {
		return err
	}

	return db.SaveDelegatorUnstakingRequest(types.NewDelegatorUnstakingRequest(staking, int64(block.Height), nodeId, delegatorID))
}