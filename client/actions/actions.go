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

type Call interface {
	// Validator
	IsAuthorizedValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (isAuthorized bool, err error)
	GetValidator(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (v structs.Validator, err error)
	GetValidatorWithInfo(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (v structs.Validator, err error)

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
	GetDelegation(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error)
	GetDelegationWithInfo(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error)
	GetDelegationState(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (ds structs.DelegationState, err error)
	GetValidatorDelegations(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (delegations []structs.Delegation, err error)
	GetHolderDelegations(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, holder common.Address) (delegations []structs.Delegation, err error)
	GetValidatorDelegationsIDs(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (delegationsIDs []uint64, err error)
}

type BCGetter interface {
	GetBoundContractCaller(ctx context.Context, addr common.Address, a abi.ABI) *bind.BoundContract
}

type Caches struct {
	Account        *lru.Cache
	AccountLock    sync.RWMutex
	Delegation     *lru.Cache
	DelegationLock sync.RWMutex
}

func NewCaches() *Caches {
	return &Caches{
		Account:    lru.New(1000),
		Delegation: lru.New(9000),
	}
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

	if ce.EventName == "RoleGranted" ||
		ce.EventName == "RoleRevoked" ||
		ce.EventName == "Upgraded" ||
		ce.EventName == "AdminChanged" {
		ce.BoundType = "none"
		return m.dataStore.SaveContractEvent(ctx, ce)
	}
	switch ce.ContractName {
	case "validator_service":
		vIDI, ok := ce.Params["validatorId"]
		if !ok {
			return errors.New("structure is not a validator, it does not have validatorID")
		}
		vID, ok := vIDI.(*big.Int)
		if !ok {
			return errors.New("structure is not a validator, it does not have validatorID")
		}

		v, err := m.c.GetValidatorWithInfo(ctx, bc, ce.BlockHeight, vID)
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

			var linkedNodes, activeNodes uint64
			for _, node := range nodes {
				node.BlockHeight = ce.BlockHeight
				if node.Status == structs.NodeStatusActive || node.Status == structs.NodeStatusLeaving {
					activeNodes++
				}
				linkedNodes++
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

			err = m.dataStore.SaveValidatorStatistic(ctx, vID, ce.BlockHeight, ce.Time, structs.ValidatorStatisticsTypeActiveNodes, new(big.Int).SetUint64(activeNodes))
			if err != nil {
				m.l.Error("error saving SaveValidatorStatistic for ValidatorStatisticsTypeActiveNodes ", zap.Error(err))
				return err
			}

			err = m.dataStore.SaveValidatorStatistic(ctx, vID, ce.BlockHeight, ce.Time, structs.ValidatorStatisticsTypeLinkedNodes, new(big.Int).SetUint64(linkedNodes))
			if err != nil {
				m.l.Error("error saving SaveValidatorStatistic for ValidatorStatisticsTypeLinkedNodes ", zap.Error(err))
				break
			}

			err = m.dataStore.UpdateCountsOfValidator(ctx, vID)
			if err != nil {
				m.l.Error("error saving SaveValidatorStatistic for UpdateNodeCountsOfValidator ", zap.Error(err))
				break
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
			return fmt.Errorf("error in nodes: %w", err)
		}

		isExitEvent := ce.EventName == "ExitCompleted"
		if isExitEvent {
			err = m.dataStore.SaveNodes(ctx, []structs.Node{n}, common.Address{})
			if err != nil {
				m.l.Error("error saving exiting/exited node", zap.Error(err))
				return err
			}
		}

		nodes, err := m.c.GetValidatorNodes(ctx, bc, ce.BlockHeight, n.ValidatorID)
		if err != nil {
			return fmt.Errorf("error getting validator nodes %w", err)
		}

		var linkedNodes, activeNodes uint64
		for _, node := range nodes {
			node.BlockHeight = ce.BlockHeight
			if node.Status == structs.NodeStatusActive || node.Status == structs.NodeStatusLeaving {
				activeNodes++
			}
			linkedNodes++
		}

		if err = m.dataStore.SaveNodes(ctx, nodes, common.Address{}); err != nil {
			return fmt.Errorf("error storing nodes %w", err)
		}

		err = m.dataStore.SaveValidatorStatistic(ctx, n.ValidatorID, ce.BlockHeight, ce.Time, structs.ValidatorStatisticsTypeActiveNodes, new(big.Int).SetUint64(activeNodes))
		if err != nil {
			m.l.Error("error saving SaveValidatorStatistic for ValidatorStatisticsTypeActiveNodes ", zap.Error(err))
			return err
		}

		if !isExitEvent {
			err = m.dataStore.SaveValidatorStatistic(ctx, n.ValidatorID, ce.BlockHeight, ce.Time, structs.ValidatorStatisticsTypeLinkedNodes, new(big.Int).SetUint64(linkedNodes))
			if err != nil {
				m.l.Error("error saving SaveValidatorStatistic for ValidatorStatisticsTypeLinkedNodes ", zap.Error(err))
				break
			}
		}

		err = m.dataStore.UpdateCountsOfValidator(ctx, n.ValidatorID)
		if err != nil {
			m.l.Error("error saving SaveValidatorStatistic for UpdateNodeCountsOfValidator ", zap.Error(err))
			break
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

		d, err := m.c.GetDelegationWithInfo(ctx, bc, ce.BlockHeight, dID)
		if err != nil {
			return fmt.Errorf("error running delegationChanged  %w", err)
		}
		d.TransactionHash = ce.TransactionHash
		d.BlockHeight = ce.BlockHeight

		m.caches.DelegationLock.Lock()
		m.caches.Delegation.Add(dID, d)
		m.caches.DelegationLock.Unlock()
		if err := m.dataStore.SaveDelegation(ctx, d); err != nil {
			return fmt.Errorf("error storing delegation %w", err)
		}

		if err := m.dataStore.SaveAccount(ctx, structs.Account{
			Address: d.Holder,
			Type:    structs.AccountTypeDelegator,
		}); err != nil {
			return fmt.Errorf("error storing account %w", err)
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
				m.caches.AccountLock.RLock()
				_, ok := m.caches.Account.Get(ad)
				m.caches.AccountLock.RUnlock()
				if !ok {
					if err := m.dataStore.SaveAccount(ctx, structs.Account{Address: ad}); err != nil {
						return err
					}
					m.caches.AccountLock.Lock()
					m.caches.Account.Add(ad, structs.Account{Address: ad})
					m.caches.AccountLock.Unlock()
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

type StateError struct {
	State  structs.DelegationState
	Err    error
	Amount *big.Int
}

type IdAmount struct {
	Id     uint64
	Amount *big.Int
}

func (m *Manager) getValidatorDelegationValues(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, blockTime time.Time, validatorID *big.Int) error {

	delegationsIDs, err := m.c.GetValidatorDelegationsIDs(ctx, bc, blockNumber, validatorID)

	if err != nil {
		if err == transport.ErrEmptyResponse {
			return nil
		}
		return fmt.Errorf("error calling GetValidatorDelegationsIDs %w", err)
	}

	contr := bc.GetContract()

	ids := make(chan IdAmount)
	defer close(ids)
	out := make(chan StateError, len(delegationsIDs))
	defer close(out)
	for i := 0; i < 10; i++ {
		go m.asyncStatuses(ctx, contr, blockNumber, ids, out)
	}

	var sent uint64
	for _, dID := range delegationsIDs {

		m.caches.DelegationLock.RLock()
		delI, ok := m.caches.Delegation.Get(dID)
		m.caches.DelegationLock.RUnlock()
		var del structs.Delegation
		if !ok {
			del, err = m.c.GetDelegation(ctx, bc, blockNumber, new(big.Int).SetUint64(dID))
			if err != nil {
				return err
			}
			m.caches.DelegationLock.Lock()
			m.caches.Delegation.Add(dID, del)
			m.caches.DelegationLock.Unlock()
		} else {
			del = delI.(structs.Delegation)
		}

		if del.State != structs.DelegationStateCOMPLETED && del.State != structs.DelegationStateCANCELED && del.State != structs.DelegationStateREJECTED {
			sent++
			ids <- IdAmount{dID, new(big.Int).Set(del.Amount)}
		}
	}

	ts := new(big.Int)
	for i := uint64(0); i < sent; i++ {
		ds := <-out
		if ds.Err != nil {
			err = ds.Err
		}
		//Total Stake
		if ds.State == structs.DelegationStateDELEGATED || ds.State == structs.DelegationStateUNDELEGATION_REQUESTED {
			ts = ts.Add(ts, ds.Amount)
		}
	}
	if err != nil {
		return err
	}

	err = m.dataStore.SaveValidatorStatistic(ctx, validatorID, blockNumber, blockTime, structs.ValidatorStatisticsTypeTotalStake, ts)
	if err != nil {
		return fmt.Errorf("error calling SaveValidatorStatistic (ValidatorStatisticsTypeTotalStake) %w", err)
	}

	err = m.dataStore.UpdateCountsOfValidator(ctx, validatorID)
	if err != nil {
		return fmt.Errorf("error calling UpdateCountsOfValidator %w", err)
	}

	return nil
}

func (m *Manager) asyncStatuses(ctx context.Context, contr *bind.BoundContract, blockNumber uint64, in chan IdAmount, out chan StateError) {
	for i := range in {
		ds, err := m.c.GetDelegationState(ctx, contr, blockNumber, new(big.Int).SetUint64(i.Id))

		m.caches.DelegationLock.Lock()
		delI, ok := m.caches.Delegation.Get(i.Id)
		if ok {
			del := delI.(structs.Delegation)
			del.State = ds
			m.caches.Delegation.Add(i.Id, del)
		}
		m.caches.DelegationLock.Unlock()

		out <- StateError{ds, err, new(big.Int).Set(i.Amount)}
	}
}
