package actions

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"go.uber.org/zap"
)

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
	var nodesInfo map[uint64]NodeAggregationInfo
	for o := range outp {
		if o.err != nil {
			errors = append(errors, o.err)
		}
		switch o.typ {
		case "validators":
			vldrs = o.data.([]structs.Validator)
		case "nodes":
			nodesInfo = o.data.(map[uint64]NodeAggregationInfo)
		}

		count--
		if count == 0 {
			break
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}

	contractCallerForDelegations := m.tr.GetBoundContractCaller(ctx, contractForDelegations.Addr, contractForDelegations.Abi)
	m.l.Info("synchronization - storing validator changes", zap.Uint64("block", currentBlock), zap.Time("blocTime", blockTime))
	for _, v := range vldrs {
		if err := m.saveValidatorStatChanges(ctx, v, currentBlock, blockTime); err != nil {
			return fmt.Errorf("error saveValidatorStatChanges %w", err)
		}

		if err := m.getValidatorDelegationValues(ctx, contractCallerForDelegations, currentBlock, blockTime, v.ValidatorID); err != nil {
			return fmt.Errorf("error getting total stake %w", err)
		}

		nInfo, ok := nodesInfo[v.ValidatorID.Uint64()]
		if ok {
			err := m.dataStore.SaveValidatorStatistic(ctx, v.ValidatorID, currentBlock, blockTime, structs.ValidatorStatisticsTypeActiveNodes, big.NewInt(int64(nInfo.ActiveNodeCount)))
			if err != nil {
				m.l.Error("error saving SaveValidatorStatistic for ValidatorStatisticsTypeActiveNodes ", zap.Error(err))
				return err
			}

			err = m.dataStore.SaveValidatorStatistic(ctx, v.ValidatorID, currentBlock, blockTime, structs.ValidatorStatisticsTypeLinkedNodes, big.NewInt(int64(nInfo.LinkedNodeCount)))
			if err != nil {
				m.l.Error("error saving SaveValidatorStatistic for ValidatorStatisticsTypeLinkedNodes ", zap.Error(err))
				return err
			}

			err = m.dataStore.UpdateCountsOfValidator(ctx, v.ValidatorID)
			if err != nil {
				m.l.Error("error saving SaveValidatorStatistic for UpdateNodeCountsOfValidator ", zap.Error(err))
				return err
			}
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
		if finished, err := m.syncDelegations(ctx, cV, *big.NewInt(i), currentBlock); finished || err != nil {
			select {
			case end <- 1:
			default:
			}
			break
		}
	}
}

func (m *Manager) syncDelegations(ctx context.Context, cV contract.ContractsContents, dID big.Int, currentBlock uint64) (finished bool, err error) {

	bc := m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi)
	var d structs.Delegation

	m.caches.DelegationLock.RLock()
	delI, ok := m.caches.Delegation.Get(&dID)
	m.caches.DelegationLock.RUnlock()
	if !ok {
		d, err = m.c.GetDelegation(ctx, bc, currentBlock, &dID)
		m.l.Debug("syncDelegations", zap.Uint64("id", dID.Uint64()), zap.Error(err))
		if err != nil {
			if err == transport.ErrEmptyResponse {
				return true, nil
			}
			m.l.Error("error occurs on sync GetDelegation", zap.Error(err))
			return true, err
		}
	} else {
		d = delI.(structs.Delegation)
	}

	d.State, err = m.c.GetDelegationState(ctx, bc.GetContract(), currentBlock, &dID)
	if err != nil {
		m.l.Error("error occurs on sync GetDelegationState", zap.Error(err))
		return true, err
	}

	d.BlockHeight = currentBlock

	if err = m.dataStore.SaveDelegation(ctx, d); err != nil {
		m.l.Error("error saving delegation ", zap.Error(err))
		return true, err
	}

	m.caches.DelegationLock.Lock()
	m.caches.Delegation.Add(&dID, d)
	m.caches.DelegationLock.Unlock()
	return false, nil
}

func (m *Manager) syncValidators(ctx context.Context, cV contract.ContractsContents, currentBlock uint64) (validators []structs.Validator, err error) {
	m.l.Info("synchronization for validator starts", zap.Uint64("block height", currentBlock))

	bc := m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi)
	vID := big.NewInt(1)
	validators = []structs.Validator{}
	var vld structs.Validator
	for err == nil {
		m.l.Debug("syncValidators", zap.Uint64("id", vID.Uint64()))
		vld, err = m.c.GetValidatorWithInfo(ctx, bc, currentBlock, vID)
		if err != nil {
			if err == transport.ErrEmptyResponse {
				m.l.Info("synchronization for validators successful.")
				return validators, nil
			}
			m.l.Error("error occurs on sync GetValidatorWithInfo", zap.Error(err))
			return validators, err
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

func (m *Manager) syncNodes(ctx context.Context, cV contract.ContractsContents, currentBlock uint64) (nodes []structs.Node, err error) {
	m.l.Info("synchronization for nodes starts", zap.Uint64("block height", currentBlock))

	bc := m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi)
	nID := big.NewInt(1)
	var n structs.Node
	nodes = []structs.Node{}
	for err == nil {
		m.l.Debug("syncNodes", zap.Uint64("id", nID.Uint64()))
		n, err = m.c.GetNodeWithInfo(ctx, bc, currentBlock, nID)
		if err != nil {
			if err == transport.ErrEmptyResponse {
				return nodes, nil
			}
			m.l.Error("error occurs on sync GetNodeWithInfo", zap.Error(err))
			return nodes, err
		}

		err = m.dataStore.SaveNodes(ctx, []structs.Node{n}, common.Address{})
		if err != nil {
			m.l.Error("error saving nodes ", zap.Error(err))
			return nodes, err
		}
		nodes = append(nodes, n)
		nID.Add(nID, big.NewInt(1))
	}

	m.l.Info("synchronization for nodes successful.")
	return nodes, nil
}

type NodeAggregationInfo struct {
	ActiveNodeCount uint64
	LinkedNodeCount uint64
}

func (m *Manager) groupNodesInfo(nodes []structs.Node) map[uint64]NodeAggregationInfo {
	nodeInfoByValidator := map[uint64]NodeAggregationInfo{}
	for _, n := range nodes {
		nInfo, ok := nodeInfoByValidator[n.ValidatorID.Uint64()]
		if !ok {
			nInfo = NodeAggregationInfo{}
		}

		nInfo.LinkedNodeCount += 1
		if n.Status == structs.NodeStatusActive {
			nInfo.ActiveNodeCount += 1
		}

		nodeInfoByValidator[n.ValidatorID.Uint64()] = nInfo
	}

	return nodeInfoByValidator
}

func (m *Manager) syncNodesAsync(ctx context.Context, cV contract.ContractsContents, currentBlock uint64, outp chan syncOutp) {
	nodes, err := m.syncNodes(ctx, cV, currentBlock)
	nodeInfoByValidator := m.groupNodesInfo(nodes)
	outp <- syncOutp{
		typ:  "nodes",
		err:  err,
		data: nodeInfoByValidator,
	}
}
