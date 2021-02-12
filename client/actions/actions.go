package actions

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/figment-networks/skale-indexer/client/standard"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"github.com/figment-networks/skale-indexer/store"

	"github.com/golang/groupcache/lru"
)

var implementedContractNames = []string{"skale_token", "delegation_controller", "validator_service", "nodes", "distributor", "punisher", "skale_manager", "bounty", "bounty_v2"}

var ErrOutOfIndex = errors.New("abi: attempting to unmarshall an empty string while arguments are expected")

type Call interface {
	// Validator
	IsAuthorizedValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (isAuthorized bool, err error)
	GetValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (v structs.Validator, err error)
	GetValidatorWithInfo(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (v structs.Validator, err error)

	// Nodes
	GetValidatorNodes(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (nodes []structs.Node, err error)
	GetNode(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, nodeID *big.Int) (n structs.Node, err error)
	GetNodeWithInfo(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, nodeID *big.Int) (n structs.Node, err error)
	GetNodeNextRewardDate(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, nodeID *big.Int) (t time.Time, err error)
	GetNodeAddress(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, nodeID *big.Int) (address common.Address, err error)

	// Distributor
	GetEarnedFeeAmountOf(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (earned, endMonth *big.Int, err error)

	// Delegation
	GetPendingDelegationsTokens(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, holderAddress common.Address) (amount *big.Int, err error)
	GetDelegation(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error)
	GetDelegationWithInfo(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error)
	GetDelegationState(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (ds structs.DelegationState, err error)
	GetValidatorDelegations(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (delegations []structs.Delegation, err error)
	GetHolderDelegations(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, holder common.Address) (delegations []structs.Delegation, err error)
}

type BCGetter interface {
	GetBoundContractCaller(ctx context.Context, addr common.Address, a abi.ABI) *bind.BoundContract
}

type Caches struct {
	Account *lru.Cache
}

func NewCaches() *Caches {
	return &Caches{Account: lru.New(1000)}
}

type Manager struct {
	dataStore store.DataStore
	c         Call
	tr        transport.EthereumTransport
	cm        *contract.Manager
	l         *zap.Logger
	caches    *Caches
}

func NewManager(c Call, dataStore store.DataStore, tr transport.EthereumTransport, cm *contract.Manager, l *zap.Logger) *Manager {
	return &Manager{
		c:         c,
		dataStore: dataStore,
		tr:        tr,
		cm:        cm,
		l:         l,
		caches:    NewCaches(),
	}
}

func (m *Manager) GetImplementedContractNames() []string {
	return implementedContractNames
}

func (m *Manager) GetBlockHeader(ctx context.Context, height big.Int) (h *types.Header, err error) {
	h, err = m.tr.GetBlockHeader(ctx, &height)
	return h, err
}

func (m *Manager) AfterEventLog(ctx context.Context, c contract.ContractsContents, ce structs.ContractEvent) (err error) {

	bc := m.tr.GetBoundContractCaller(ctx, c.Addr, c.Abi)

	if ce.EventName == "RoleGranted" {
		// BUG(lukanus): save this in correct form
		/*
			"inputs": [
				{
					"indexed": true,
					"internalType": "bytes32",
					"name": "role",
					"type": "bytes32"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "account",
					"type": "address"
				},
				{
					"indexed": true,
					"internalType": "address",
					"name": "sender",
					"type": "address"
				}
			],*/
		ce.BoundType = "none"
		return m.dataStore.SaveContractEvent(ctx, ce)
	}

	switch ce.ContractName {
	case "validator_service":
		vIDI, ok := ce.Params["validatorId"]
		if !ok {
			return errors.New("structure is not a validator, it does not have valiadtorId")
		}
		vID, ok := vIDI.(*big.Int)
		if !ok {
			return errors.New("structure is not a validator, it does not have valiadtorId")
		}

		v, err := m.c.GetValidatorWithInfo(ctx, bc.GetContract(), ce.BlockHeight, vID)
		if err != nil {
			return fmt.Errorf("error running validatorChanged  %w", err)
		}
		v.BlockHeight = ce.BlockHeight
		if err = m.dataStore.SaveValidator(ctx, v); err != nil {
			return fmt.Errorf("error storing validator %w", err)
		}

		if err = m.saveValidatorStatChanges(ctx, v, ce.BlockHeight, ce.Time); err != nil {
			return fmt.Errorf("error storing changes %w", err)
		}

		if ce.EventName == "NodeAddressWasAdded" || ce.EventName == "NodeAddressWasRemoved" {
			cV, ok := m.cm.GetContractByNameVersion("nodes", c.Version)
			if !ok {
				return errors.New("Node contract is not found for version :" + c.Version)
			}
			nodes, err := m.c.GetValidatorNodes(ctx, m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi), ce.BlockHeight, vID)
			if err != nil {
				return fmt.Errorf("error getting validator nodes %w", err)
			}

			for _, node := range nodes {
				node.BlockHeight = ce.BlockHeight
			}

			removedNodeAddr := common.Address{}
			if ce.EventName == "NodeAddressWasRemoved" {
				nodeAddrI, ok := ce.Params["nodeAddress"]
				if !ok {
					return errors.New("structure is not for NodeAddressWasAdded or NodeAddressWasRemoved, it does not have nodeAddress")
				}
				removedNodeAddr, ok = nodeAddrI.(common.Address)
				if !ok {
					return errors.New("structure is not for NodeAddressWasAdded or NodeAddressWasRemoved, it does not have nodeAddress")
				}
			}

			if err := m.dataStore.SaveNodes(ctx, nodes, removedNodeAddr); err != nil {
				return fmt.Errorf("error storing validator nodes %w", err)
			}

			vsp := structs.ValidatorStatisticsParams{
				ValidatorID: vID.String(),
				BlockHeight: ce.BlockHeight,
				BlockTime:   ce.Time,
			}

			if err := m.dataStore.CalculateActiveNodes(ctx, vsp); err != nil {
				return fmt.Errorf("error calculating active nodes %w", err)
			}

			if err := m.dataStore.CalculateLinkedNodes(ctx, vsp); err != nil {
				return fmt.Errorf("error calculating linked nodes %w", err)
			}

		} else if ce.EventName == "ValidatorRegistered" {

			if err = m.dataStore.SaveAccount(ctx, structs.Account{
				Address: v.ValidatorAddress,
				Type:    structs.AccountTypeValidator,
			}); err != nil {
				return fmt.Errorf("error storing account %w", err)
			}

			if err = m.dataStore.SaveSystemEvent(ctx, structs.SystemEvent{
				Height: ce.BlockHeight,
				Time:   ce.Time,
				Kind:   structs.SysEvtTypeFeeChanged,
				After:  *v.FeeRate,
			}); err != nil {
				return fmt.Errorf("error storing system event %w", err)
			}

			if err = m.dataStore.SaveSystemEvent(ctx, structs.SystemEvent{
				Height: ce.BlockHeight,
				Time:   ce.Time,
				Kind:   structs.SysEvtTypeMDRChanged,
				After:  *v.MinimumDelegationAmount,
			}); err != nil {
				return fmt.Errorf("error storing system event %w", err)
			}

		} else if ce.EventName == "ValidatorAddressChanged" {
			newAddrI, ok := ce.Params["newAddress"]
			if !ok {
				return errors.New("structure is not for ValidatorAddressChanged, it does not have newAddress")
			}
			addr, ok := newAddrI.(common.Address)
			if !ok {
				return errors.New("structure is not for ValidatorAddressChanged, it does not have newAddress")
			}

			if err := m.dataStore.SaveAccount(ctx, structs.Account{
				Address: addr,
				Type:    structs.AccountTypeValidator,
			}); err != nil {
				return fmt.Errorf("error storing account %w", err)
			}
		} else if ce.EventName == "ValidatorWasEnabled" {
			if err := m.dataStore.SaveSystemEvent(ctx, structs.SystemEvent{
				Height:      ce.BlockHeight,
				Time:        ce.Time,
				Kind:        structs.SysEvtTypeJoinedActiveSet,
				RecipientID: *vID}); err != nil {
				return fmt.Errorf("error storing system event %w", err)
			}
		} else if ce.EventName == "ValidatorWasDisabled" {
			if err := m.dataStore.SaveSystemEvent(ctx, structs.SystemEvent{
				Height:      ce.BlockHeight,
				Time:        ce.Time,
				Kind:        structs.SysEvtTypeLeftActiveSet,
				RecipientID: *vID}); err != nil {
				return fmt.Errorf("error storing system event %w", err)
			}
		}

		ce.BoundType = "validator"
		ce.BoundID = append(ce.BoundID, *vID)
	case "nodes":

		nIDI, ok := ce.Params["nodeIndex"]
		if !ok {
			return errors.New("structure is not a node, it does not have nodeIndex")
		}
		nID, ok := nIDI.(*big.Int)
		if !ok {
			return errors.New("structure is not a node, it does not have nodeIndex")
		}

		n, err := m.c.GetNode(ctx, bc, ce.BlockHeight, nID)
		if err != nil {
			// TODO: change err message from line 203
			return errors.New("structure is not a node")
		}
		n.BlockHeight = ce.BlockHeight
		if err = m.dataStore.SaveNodes(ctx, []structs.Node{n}, common.Address{}); err != nil {
			return fmt.Errorf("error storing nodes %w", err)
		}
		vs := structs.ValidatorStatisticsParams{
			ValidatorID: n.ValidatorID.String(),
			BlockHeight: ce.BlockHeight,
			BlockTime:   ce.Time,
		}
		err = m.dataStore.CalculateActiveNodes(ctx, vs)
		if err != nil {
			return fmt.Errorf("error calculating active nodes %w", err)
		}
		err = m.dataStore.CalculateLinkedNodes(ctx, vs)
		if err != nil {
			return fmt.Errorf("error calculating linked nodes %w", err)
		}
		ce.BoundType = "node"
		ce.BoundID = append(ce.BoundID, *nID)

	case "punisher":

		switch ce.EventName {
		case "slash":
			vIDI, ok := ce.Params["validatorId"]
			if !ok {
				return errors.New("structure is not a slash, it does not have validatorId")
			}
			vID, ok := vIDI.(*big.Int)
			if !ok {
				return errors.New("structure is not a slash, it does not have validatorId")
			}

			amountI, ok := ce.Params["amount"]
			if !ok {
				return errors.New("structure is not a slash")
			}

			am, ok := amountI.(*big.Int)
			if !ok {
				return errors.New("structure is not a slash")
			}

			if err = m.dataStore.SaveSystemEvent(ctx, structs.SystemEvent{
				Height: ce.BlockHeight,
				Time:   ce.Time,
				Kind:   structs.SysEvtTypeSlashed,
				After:  *am,
			}); err != nil {
				return fmt.Errorf("error storing system event %w", err)
			}

			ce.BoundType = "validator"
			ce.BoundID = append(ce.BoundID, *vID)
		case "forgive":
			wAddrI, ok := ce.Params["wallet"]
			if !ok {
				return errors.New("structure is not a forgive, it does not have wallet")
			}
			wAddr, ok := wAddrI.(common.Address)
			if !ok {
				return errors.New("structure is not a forgive")
			}

			amountI, ok := ce.Params["amount"]
			if !ok {
				return errors.New("structure is not a forgive")
			}

			am, ok := amountI.(*big.Int)
			if !ok {
				return errors.New("structure is not a forgive")
			}

			if err = m.dataStore.SaveSystemEvent(ctx, structs.SystemEvent{
				Height:    ce.BlockHeight,
				Time:      ce.Time,
				Kind:      structs.SysEvtTypeForgiven,
				Recipient: wAddr,
				After:     *am,
			}); err != nil {
				return fmt.Errorf("error storing system forgive %w", err)
			}

			ce.BoundAddress = append(ce.BoundAddress, wAddr)
		}

	case "distributor":
		switch ce.EventName {
		case "WithdrawBounty":
			hAddrI, ok := ce.Params["holder"]
			if !ok {
				return errors.New("structure is not a distributor, it does not have holder")
			}
			hAddr, ok := hAddrI.(common.Address)
			if !ok {
				return errors.New("structure is not a distributor, it does not have holder")
			}
			ce.BoundAddress = append(ce.BoundAddress, hAddr)
			dAddrI, ok := ce.Params["destination"]
			if !ok {
				return errors.New("structure is not a distributor, it does not have destination")
			}
			dAddr, ok := dAddrI.(common.Address)
			if !ok {
				return errors.New("structure is not a distributor, it does not have destination")
			}
			ce.BoundAddress = append(ce.BoundAddress, dAddr)
		case "WithdrawFee":
			dAddrI, ok := ce.Params["destination"]
			if !ok {
				return errors.New("structure is not a distributor, it does not have destination")
			}
			dAddr, ok := dAddrI.(common.Address)
			if !ok {
				return errors.New("structure is not a distributor, it does not have destination")
			}
			ce.BoundAddress = append(ce.BoundAddress, dAddr)
		}

		vIDI, ok := ce.Params["validatorId"]
		if !ok {
			return errors.New("structure is not a distributor, it does not have validatorId")
		}
		vID, ok := vIDI.(*big.Int)
		if !ok {
			return errors.New("structure is not a distributor, it does not have validatorId")
		}

		ce.BoundType = "validator"
		ce.BoundID = append(ce.BoundID, *vID)

	case "delegation_controller":

		dIDI, ok := ce.Params["delegationId"]
		if !ok {
			return errors.New("structure is not a delegation, it does not have delegationId")
		}
		dID, ok := dIDI.(*big.Int)
		if !ok {
			return errors.New("structure is not a delegation, it does not have delegationId")
		}

		d, err := m.c.GetDelegationWithInfo(ctx, bc.GetContract(), ce.BlockHeight, dID)
		if err != nil {
			return fmt.Errorf("error running delegationChanged  %w", err)
		}
		d.TransactionHash = ce.TransactionHash
		d.BlockHeight = ce.BlockHeight

		if err := m.dataStore.SaveDelegation(ctx, d); err != nil {
			return fmt.Errorf("error storing delegation %w", err)
		}

		if err := m.dataStore.SaveAccount(ctx, structs.Account{
			Address: d.Holder,
			Type:    structs.AccountTypeDelegator,
		}); err != nil {
			return fmt.Errorf("error storing account %w", err)
		}

		if err := m.dataStore.CalculateTotalStake(ctx, structs.ValidatorStatisticsParams{
			ValidatorID: d.ValidatorID.String(),
			BlockHeight: ce.BlockHeight,
			BlockTime:   ce.Time,
		}); err != nil {
			return fmt.Errorf("error calculating total stake %w", err)
		}

		sysEvt := structs.SystemEvent{
			Height: ce.BlockHeight,
			Time:   ce.Time,
			After:  *d.Amount,
		}
		switch ce.EventName {
		case "DelegationProposed":
			sysEvt.Kind = structs.SysEvtTypeNewDelegation
			sysEvt.Sender = d.Holder
			sysEvt.RecipientID = *d.ValidatorID
		case "DelegationAccepted":
			sysEvt.Kind = structs.SysEvtTypeDelegationAccepted
			sysEvt.Recipient = d.Holder
			sysEvt.SenderID = *d.ValidatorID
		case "DelegationRequestCanceledByUser":
			sysEvt.Kind = structs.SysEvtTypeDelegationRejected
			sysEvt.Sender = d.Holder
			sysEvt.RecipientID = *d.ValidatorID
		case "UndelegationRequested":
			sysEvt.Kind = structs.SysEvtTypeUndeledationRequested
			sysEvt.Sender = d.Holder
			sysEvt.RecipientID = *d.ValidatorID
		}

		if err = m.dataStore.SaveSystemEvent(ctx, sysEvt); err != nil {
			return fmt.Errorf("error storing system event %w", err)
		}

		ce.BoundType = "delegation"
		ce.BoundID = []big.Int{*dID, *d.ValidatorID}
	case "skale_manager":

		nIDI, ok := ce.Params["nodeIndex"]
		if !ok {
			return errors.New("structure is not a skale_manager, it does not have nodeIndex")
		}
		nID, ok := nIDI.(*big.Int)
		if !ok {
			return errors.New("structure is not a skale_manager, it does not have nodeIndex")
		}

		cV, ok := m.cm.GetContractByNameVersion("nodes", c.Version)
		if !ok {
			return errors.New("Node contract is not found for version :" + c.Version)
		}

		n, err := m.c.GetNodeWithInfo(ctx, m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi), ce.BlockHeight, nID)
		if err != nil {
			return errors.New("structure is not a skale_manager")
		}

		n.BlockHeight = ce.BlockHeight
		if err = m.dataStore.SaveNodes(ctx, []structs.Node{n}, common.Address{}); err != nil {
			return fmt.Errorf("error storing node %w", err)
		}

		ce.BoundID = append(ce.BoundID, *nID)
		ce.BoundType = "node"
	case "skale_token":
		if ce.EventName == "Transfer" || ce.EventName == "Approval" {
			ce, err = standard.DecodeERC20Events(ctx, ce)
			if err != nil {
				return fmt.Errorf("error decoding event ERC20 %w", err)
			}
			for _, ad := range ce.BoundAddress {
				if _, ok := m.caches.Account.Get(ad); !ok {
					if err := m.dataStore.SaveAccount(ctx, structs.Account{Address: ad}); err != nil {
						return err
					}
					m.caches.Account.Add(ad, structs.Account{Address: ad})
				}
			}
		}
		ce.BoundType = "token"

	default:
		m.l.Debug("Unknown event type", zap.String("type", ce.ContractName), zap.Any("event", ce))
	}

	return m.dataStore.SaveContractEvent(ctx, ce)

}

func (m *Manager) saveValidatorStatChanges(ctx context.Context, validator structs.Validator, blockNumber uint64, blockTime time.Time) error {

	err := m.dataStore.SaveValidatorStatistic(ctx, validator.ValidatorID, blockNumber, blockTime, structs.ValidatorStatisticsTypeFee, validator.FeeRate)
	if err != nil {
		return fmt.Errorf("error calling SaveValidatorStatistic (ValidatorStatisticsTypeFee) %w", err)
	}

	err = m.dataStore.SaveValidatorStatistic(ctx, validator.ValidatorID, blockNumber, blockTime, structs.ValidatorStatisticsTypeMDR, validator.MinimumDelegationAmount)
	if err != nil {
		return fmt.Errorf("error calling SaveValidatorStatistic (ValidatorStatisticsTypeMDR) %w", err)
	}

	err = m.dataStore.SaveValidatorStatistic(ctx, validator.ValidatorID, blockNumber, blockTime, structs.ValidatorStatisticsTypeAuthorized, boolToBigInt(validator.Authorized))
	if err != nil {
		return fmt.Errorf("error calling SaveValidatorStatistic (ValidatorStatisticsTypeAuthorized) %w", err)
	}

	if validator.ValidatorAddress.String() != "" {
		err = m.dataStore.SaveValidatorStatistic(ctx, validator.ValidatorID, blockNumber, blockTime, structs.ValidatorStatisticsTypeValidatorAddress, validator.ValidatorAddress.Hash().Big())
		if err != nil {
			return fmt.Errorf("error calling SaveValidatorStatistic (ValidatorStatisticsTypeAuthorized) %w", err)
		}
	}

	if validator.RequestedAddress.String() != "" {
		err = m.dataStore.SaveValidatorStatistic(ctx, validator.ValidatorID, blockNumber, blockTime, structs.ValidatorStatisticsTypeRequestedAddress, validator.RequestedAddress.Hash().Big())
		if err != nil {
			return fmt.Errorf("error calling SaveValidatorStatistic (ValidatorStatisticsTypeAuthorized) %w", err)
		}
	}

	return nil
}

func boolToBigInt(a bool) *big.Int {
	if a {
		return big.NewInt(1)
	}
	return big.NewInt(0)
}

type syncOutp struct {
	typ  string
	err  error
	data interface{}
}

func (m *Manager) SyncForBeginningOfEpoch(ctx context.Context, version string, currentBlock uint64, blockTime time.Time) error {
	m.l.Info("synchronization starts", zap.Uint64("block", currentBlock), zap.Time("blockTime", blockTime))

	contractForValidator, ok := m.cm.GetContractByNameVersion("validator_service", version)
	if !ok {
		m.l.Error("failed to synchronize validators. contract is not found.")
		return errors.New("contract is not found for validators for version :" + version)
	}

	contractForNodes, ok := m.cm.GetContractByNameVersion("nodes", version)
	if !ok {
		m.l.Error("failed to synchronize nodes. contract is not found.")
		return errors.New("contract is not found for nodes for version :" + version)
	}

	contractForDelegations, ok := m.cm.GetContractByNameVersion("delegation_controller", version)
	if !ok {
		m.l.Error("failed to synchronize delegations. contract is not found.")
		return errors.New("contract is not found for delegations for version :" + version)
	}

	outp := make(chan syncOutp, 3)
	defer close(outp)
	go m.syncValidatorsAsync(ctx, contractForValidator, currentBlock, outp)
	go m.syncNodesAsync(ctx, contractForNodes, currentBlock, outp)
	go m.syncDelegationsAsync(ctx, contractForDelegations, currentBlock, outp)

	var count = 3
	var errors []error
	var vldrs []structs.Validator
	for o := range outp {
		if o.err != nil {
			errors = append(errors, o.err)
		}
		if o.typ == "validators" {
			vldrs = o.data.([]structs.Validator)
		}

		count--
		if count == 0 {
			break
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}

	m.l.Info("synchronization - storing validator changes", zap.Uint64("block", currentBlock), zap.Time("blocTime", blockTime))
	for _, v := range vldrs {
		if err := m.saveValidatorStatChanges(ctx, v, currentBlock, blockTime); err != nil {
			m.l.Error(err.Error())
		}
		vs := structs.ValidatorStatisticsParams{
			ValidatorID: v.ValidatorID.String(),
			BlockHeight: currentBlock,
			BlockTime:   blockTime,
		}
		if err := m.dataStore.CalculateTotalStake(ctx, vs); err != nil {
			m.l.Error(err.Error())
		}
		if err := m.dataStore.CalculateActiveNodes(ctx, vs); err != nil {
			m.l.Error(err.Error())
		}
		if err := m.dataStore.CalculateLinkedNodes(ctx, vs); err != nil {
			m.l.Error(err.Error())
		}
	}

	m.l.Info("synchronization successfully finishes", zap.Uint64("block", currentBlock), zap.Time("blocTime", blockTime))

	return nil
}

func populate(ch, end chan int64) {
	var i int64
	for {
		select {
		case <-end:
			close(ch)
			return
		case ch <- i:
			i++
		}
	}
}

func (m *Manager) syncDelegationsAsync(ctx context.Context, cV contract.ContractsContents, currentBlock uint64, outp chan syncOutp) {
	m.l.Info("synchronization for delegations starts", zap.Uint64("block height", currentBlock))
	wg := &sync.WaitGroup{}
	ch := make(chan int64)
	end := make(chan int64)
	defer close(end)
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go m.syncDelegationsAsyncC(ctx, cV, currentBlock, ch, end, wg)
	}
	go populate(ch, end)
	wg.Wait()

	m.l.Info("sending delegations")
	outp <- syncOutp{
		typ: "delegations",
	}

	m.l.Info("synchronization for delegations successful.")
}

func (m *Manager) syncDelegationsAsyncC(ctx context.Context, cV contract.ContractsContents, currentBlock uint64, in, end chan int64, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range in {
		if err := m.syncDelegations(ctx, cV, *big.NewInt(i), currentBlock); err != nil {
			m.l.Debug("populating end", zap.Error(err))
			select {
			case end <- 1:
			default:
			}

			break
		}
	}
}

func (m *Manager) syncDelegations(ctx context.Context, cV contract.ContractsContents, dID big.Int, currentBlock uint64) (err error) {

	bc := m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi)
	var d structs.Delegation

	d, err = m.c.GetDelegationWithInfo(ctx, bc.GetContract(), currentBlock, &dID)
	m.l.Debug("syncDelegations", zap.Uint64("id", dID.Uint64()), zap.Error(err))
	if err != nil {
		if err.Error() != ErrOutOfIndex.Error() {
			m.l.Error("error occurs on sync GetDelegationWithInfo", zap.Error(err))
		}
		return err
	}
	d.BlockHeight = currentBlock
	err = m.dataStore.SaveDelegation(ctx, d)
	if err != nil {
		m.l.Error("error saving delegation ", zap.Error(err))
		return err
	}

	return nil
}

func (m *Manager) syncValidators(ctx context.Context, cV contract.ContractsContents, currentBlock uint64) (validators []structs.Validator, err error) {
	m.l.Info("synchronization for validator starts", zap.Uint64("block height", currentBlock))

	bc := m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi)
	vID := big.NewInt(1)
	validators = []structs.Validator{}
	var vld structs.Validator
	for err == nil {
		m.l.Debug("syncValidators", zap.Uint64("id", vID.Uint64()))
		vld, err = m.c.GetValidatorWithInfo(ctx, bc.GetContract(), currentBlock, vID)
		if err != nil {
			if err.Error() != ErrOutOfIndex.Error() {
				m.l.Error("error occurs on sync GetValidatorWithInfo", zap.Error(err))
				return validators, err
			}
			continue
		}
		vld.BlockHeight = currentBlock
		err = m.dataStore.SaveValidator(ctx, vld)
		if err != nil {
			m.l.Error("error saving validators ", zap.Error(err))
			return validators, err
		}
		vID.Add(vID, big.NewInt(1))
		validators = append(validators, vld)
	}

	m.l.Info("synchronization for validators successful.")
	return validators, nil
}

func (m *Manager) syncValidatorsAsync(ctx context.Context, cV contract.ContractsContents, currentBlock uint64, outp chan syncOutp) {
	v, err := m.syncValidators(ctx, cV, currentBlock)
	outp <- syncOutp{
		typ:  "validators",
		err:  err,
		data: v,
	}
}

func (m *Manager) syncNodes(ctx context.Context, cV contract.ContractsContents, currentBlock uint64) (err error) {
	m.l.Info("synchronization for nodes starts", zap.Uint64("block height", currentBlock))

	bc := m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi)
	nID := big.NewInt(1)
	var n structs.Node
	for err == nil {
		m.l.Debug("syncNodes", zap.Uint64("id", nID.Uint64()))
		n, err = m.c.GetNodeWithInfo(ctx, bc, currentBlock, nID)
		if err != nil {
			if err.Error() != ErrOutOfIndex.Error() {
				m.l.Error("error occurs on sync GetNodeWithInfo", zap.Error(err))
				return err
			}
			continue
		}
		err = m.dataStore.SaveNodes(ctx, []structs.Node{n}, common.Address{})
		if err != nil {
			m.l.Error("error saving nodes ", zap.Error(err))
			return err
		}
		nID.Add(nID, big.NewInt(1))
	}

	m.l.Info("synchronization for nodes successful.")
	return nil
}

func (m *Manager) syncNodesAsync(ctx context.Context, cV contract.ContractsContents, currentBlock uint64, outp chan syncOutp) {
	err := m.syncNodes(ctx, cV, currentBlock)
	outp <- syncOutp{
		typ: "nodes",
		err: err,
	}
}
