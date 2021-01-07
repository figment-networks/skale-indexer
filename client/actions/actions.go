package actions

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
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
)

var implementedContractNames = []string{"skale_token", "delegation_controller", "validator_service", "nodes", "distributor", "punisher", "skale_manager", "bounty", "bounty_v2"}

type Call interface {
	// Validator
	IsAuthorizedValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (isAuthorized bool, err error)
	GetValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (v structs.Validator, err error)

	// Nodes
	GetValidatorNodes(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (nodes []structs.Node, err error)
	GetNode(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (n structs.Node, err error)
	GetNodeNextRewardDate(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (t time.Time, err error)

	// Distributor
	GetEarnedFeeAmountOf(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (earned, endMonth *big.Int, err error)

	// Delegation
	GetPendingDelegationsTokens(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, holderAddress common.Address) (amount *big.Int, err error)
	GetDelegation(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error)
	GetDelegationState(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (ds structs.DelegationState, err error)
	GetValidatorDelegations(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (delegations []structs.Delegation, err error)
	GetHolderDelegations(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, holder common.Address) (delegations []structs.Delegation, err error)
}

type BCGetter interface {
	GetBoundContractCaller(ctx context.Context, addr common.Address, a abi.ABI) *bind.BoundContract
}

type Manager struct {
	dataStore store.DataStore
	c         Call
	tr        transport.EthereumTransport
	cm        *contract.Manager
	l         *zap.Logger
}

func NewManager(c Call, dataStore store.DataStore, tr transport.EthereumTransport, cm *contract.Manager, l *zap.Logger) *Manager {
	return &Manager{c: c, dataStore: dataStore, tr: tr, cm: cm, l: l}
}

func (m *Manager) GetImplementedContractNames() []string {
	return implementedContractNames
}

func (m *Manager) GetBlockHeader(ctx context.Context, height *big.Int) (h *types.Header, err error) {
	h, err = m.tr.GetBlockHeader(ctx, height)
	return h, err
}

func (m *Manager) AfterEventLog(ctx context.Context, c contract.ContractsContents, ce structs.ContractEvent) (err error) {

	bc := m.tr.GetBoundContractCaller(ctx, c.Addr, c.Abi)

	switch ce.ContractName {
	case "validator_service":
		/*
			@dev Emitted when a validator registers.
				event ValidatorRegistered(
				uint validatorId
			);
			@dev Emitted when a validator address changes.
			event ValidatorAddressChanged(
				uint validatorId,
				address newAddress
			);
			@dev Emitted when a validator is enabled.
			event ValidatorWasEnabled(
				uint validatorId
			);
			@dev Emitted when a validator is disabled.
			event ValidatorWasDisabled(
				uint validatorId
			);
			@dev Emitted when a node address is linked to a validator.
			event NodeAddressWasAdded(
				uint validatorId,
				address nodeAddress
			);
			@dev Emitted when a node address is unlinked from a validator.
			event NodeAddressWasRemoved(
				uint validatorId,
				address nodeAddress
			);
		*/

		vIDI, ok := ce.Params["validatorId"]
		if !ok {
			return errors.New("Structure is not a validator")
		}

		vID, ok := vIDI.(*big.Int)
		if !ok {
			return errors.New("Structure is not a validator")
		}

		v, err := m.getValidatorChanged(ctx, bc, ce.BlockHeight, vID)
		if err != nil {
			return fmt.Errorf("error running validatorChanged  %w", err)
		}
		v.BlockHeight = ce.BlockHeight
		v.RegistrationTime = ce.Time
		//  BUG(lukanus): error storing validator sql: converting argument $1 type: unsupported type big.Int, a struct
		if err = m.dataStore.SaveValidator(ctx, v); err != nil {
			return fmt.Errorf("error storing validator %w", err)
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

			// TODO: batch insert pq: invalid byte sequence for encoding \"UTF8\": 0x00"
			for _, node := range nodes {
				node.EventTime = ce.Time
				if err := m.dataStore.SaveNode(ctx, node); err != nil {
					return fmt.Errorf("error storing validator nodes %w", err)
				}
			}
			qqq := structs.ValidatorStatisticsParams{
				ValidatorId: vID.String(),
				BlockHeight: ce.BlockHeight,
			}
			err = m.dataStore.CalculateActiveNodes(ctx, qqq)
			if err != nil {
				return fmt.Errorf("error calculating active nodes %w", err)
			}
			err = m.dataStore.CalculateLinkedNodes(ctx, qqq)
			if err != nil {
				return fmt.Errorf("error calculating linked nodes %w", err)
			}
		} else if ce.EventName == "ValidatorRegistered" {
			err = m.dataStore.SaveAccount(ctx, structs.Account{
				Address:     v.ValidatorAddress,
				AccountType: "validator",
			})
			if err != nil {
				return fmt.Errorf("error storing account %w", err)
			}
		} else if ce.EventName == "ValidatorAddressChanged" {
			newAddrI, ok := ce.Params["newAddress"]
			if !ok {
				return errors.New("Structure is not a validator")
			}

			addr, ok := newAddrI.(common.Address)
			if !ok {
				return errors.New("Structure is not a validator")
			}

			err = m.dataStore.SaveAccount(ctx, structs.Account{
				Address:     addr,
				AccountType: "validator",
			})
			if err != nil {
				return fmt.Errorf("error storing account %w", err)
			}
		}

		ce.BoundType = "validator"
		ce.BoundID = append(ce.BoundID, *vID)
	case "nodes":
		/*
			@dev Emitted when a node is created.
			event NodeCreated(
				uint nodeIndex,
				address owner,
				string name,
				bytes4 ip,
				bytes4 publicIP,
				uint16 port,
				uint16 nonce,
				uint time,
				uint gasSpend
			);

			@dev Emitted when a node completes a network exit.
			event ExitCompleted(
				uint nodeIndex,
				uint time,
				uint gasSpend
			);

			@dev Emitted when a node begins to exit from the network.
			event ExitInitialized(
				uint nodeIndex,
				uint startLeavingPeriod,
				uint time,
				uint gasSpend
			);
		*/

		nIDI, ok := ce.Params["nodeIndex"]
		if !ok {
			return errors.New("Structure is not a node")
		}
		nID, ok := nIDI.(*big.Int)
		if !ok {
			return errors.New("Structure is not a validator")
		}

		n, err := m.c.GetNode(ctx, bc, ce.BlockHeight, nID)
		if err != nil {
			// TODO: change err message from line 203
			return errors.New("Structure is not a node")
		}

		n.EventTime = ce.Time
		if err = m.dataStore.SaveNode(ctx, n); err != nil {
			return fmt.Errorf("error storing nodes %w", err)
		}
		vs := structs.ValidatorStatisticsParams{
			ValidatorId: n.ValidatorID.String(),
			BlockHeight: ce.BlockHeight,
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

		// TODO(lukanus): Get Validator Nodes maybe?
	case "punisher":
		/*
			@dev Emitted upon slashing condition.
			event Slash(
				uint validatorId,
				uint amount
			);
			@dev Emitted upon forgive condition.
			event Forgive(
				address wallet,
				uint amount
			);
		*/
		switch ce.EventName {
		case "slash":
			vIDI, ok := ce.Params["validatorId"]
			if !ok {
				return errors.New("Structure is not a validator")
			}
			vID, ok := vIDI.(*big.Int)
			if !ok {
				return errors.New("structure is not a validator")
			}
			amntI, ok := ce.Params["amount"]
			if !ok {
				return errors.New("structure is not an amount")
			}
			//amnt, ok := amntI.(*big.Int)
			_, ok = amntI.(*big.Int)
			if !ok {
				return errors.New("structure is not an amount")
			}
			/*
				var t = int(structs.ValidatorStatisticsTypeUnclaimedRewards)
				params := structs.ValidatorStatisticsParams{
					ValidatorId:      vID.String(),
					StatisticsTypeVS: string(t),
				}
				vldStatistics, err := m.dataStore.GetValidatorStatistics(ctx, params)
				if err != nil {
					return fmt.Errorf("error getting validator statistics %w", err)
				}
				if vldStatistics != nil {
					amnt.Sub(vldStatistics[0].Amount, amnt)
				}
				err = m.updateUnclaimedRewardsValidator(ctx, vID, amnt, ce.BlockHeight)
				if err != nil {
					return fmt.Errorf("error updating unclaimed rewards for validator function %w", err)
				}
			*/
			ce.BoundType = "validator"
			ce.BoundID = append(ce.BoundID, *vID)
		case "forgive":
			wAddrI, ok := ce.Params["wallet"]
			if !ok {
				return errors.New("structure is not a validator")
			}

			wAddr, ok := wAddrI.(common.Address)
			if !ok {
				return errors.New("structure is not a validator")
			}

			ce.BoundAddress = append(ce.BoundAddress, wAddr)
		}

	case "distributor":
		/*
			@dev Emitted when bounty is withdrawn.
			event WithdrawBounty(
				address holder,
				uint validatorId,
				address destination,
				uint amount
			);
			@dev Emitted when a validator fee is withdrawn.
			event WithdrawFee(
				uint validatorId,
				address destination,
				uint amount
			);
			@dev Emitted when bounty is distributed.
			event BountyWasPaid(
				uint validatorId,
				uint amount
			);
		*/

		if ce.EventName == "WithdrawFee" {
			// for unclaimed validator
			vID, ok := ce.Params["validatorId"]
			if !ok {
				return errors.New("structure is not a validator")
			}
			earned, _, err := m.c.GetEarnedFeeAmountOf(ctx, bc, ce.BlockHeight, vID.(*big.Int))
			if err != nil {
				return fmt.Errorf("error calling getEarnedFeeAmountOf function %w", err)
			}
			err = m.updateUnclaimedRewardsValidator(ctx, vID.(*big.Int), earned, ce.BlockHeight)
			if err != nil {
				return fmt.Errorf("error updating unclaimed rewards for validator function %w", err)
			}

			// for claimed validator
			amntI, ok := ce.Params["amount"]
			if !ok {
				return errors.New("structure is not an amount")
			}
			amnt, ok := amntI.(*big.Int)
			if !ok {
				return errors.New("structure is not an amount")
			}
			var t = int(structs.ValidatorStatisticsTypeClaimedRewards)
			params := structs.ValidatorStatisticsParams{
				ValidatorId:      vID.(*big.Int).String(),
				StatisticsTypeVS: string(t),
			}
			vldStatistics, err := m.dataStore.GetValidatorStatistics(ctx, params)
			if err != nil {
				return fmt.Errorf("error getting validator statistics %w", err)
			}
			if vldStatistics != nil {
				amnt.Add(amnt, vldStatistics[0].Amount)
			}
			err = m.updateClaimedRewardsValidator(ctx, vID.(*big.Int), amnt, ce.BlockHeight)
			if err != nil {
				return fmt.Errorf("error updating claimed rewards for validator function %w", err)
			}

		} else if ce.EventName == "WithdrawBounty" {
			// update unclaimed delegator
			holderI, ok := ce.Params["holder"]
			if !ok {
				return errors.New("structure is not a holder")
			}
			holder, ok := holderI.(common.Address)
			if !ok {
				return errors.New("structure is not a holder")
			}
			amntI, ok := ce.Params["amount"]
			if !ok {
				return errors.New("structure is not an amount")
			}
			amnt, ok := amntI.(*big.Int)
			if !ok {
				return errors.New("structure is not an amount")
			}
			var t = int(structs.DelegatorStatisticsTypeUnclaimedRewards)
			params := structs.DelegatorStatisticsParams{
				Holder:           holder.String(),
				StatisticsTypeDS: strconv.Itoa(t),
			}
			dlgStatistics, err := m.dataStore.GetDelegatorStatistics(ctx, params)
			if err != nil {
				return fmt.Errorf("error getting delegator statistics %w", err)
			}
			unclRwrd := big.NewInt(0)
			unclRwrd.Add(amnt, big.NewInt(0))
			if dlgStatistics != nil {
				unclRwrd.Sub(dlgStatistics[0].Amount, unclRwrd)
			}
			err = m.updateUnclaimedRewardsDelegator(ctx, holder, unclRwrd, ce.BlockHeight)
			if err != nil {
				return fmt.Errorf("error updating unclaimed delegator statistics %w", err)
			}

			t = int(structs.DelegatorStatisticsTypeClaimedRewards)
			params = structs.DelegatorStatisticsParams{
				Holder:           holder.String(),
				StatisticsTypeDS: strconv.Itoa(t),
			}
			dlgStatistics, err = m.dataStore.GetDelegatorStatistics(ctx, params)
			if err != nil {
				return fmt.Errorf("error getting delegator statistics %w", err)
			}
			if dlgStatistics != nil {
				amnt.Add(amnt, dlgStatistics[0].Amount)
			}
			err = m.updateClaimedRewardsDelegator(ctx, holder, amnt, ce.BlockHeight)
			if err != nil {
				return fmt.Errorf("error getting delegator statistics %w", err)
			}
		} else if ce.EventName == "BountyWasPaid" {
			// update unclaimed rewards validator
			vIDI, ok := ce.Params["validatorId"]
			if !ok {
				return errors.New("Structure is not a validator")
			}
			vID, ok := vIDI.(*big.Int)
			if !ok {
				return errors.New("Structure is not a validator")
			}
			amntI, ok := ce.Params["amount"]
			if !ok {
				return errors.New("Structure is not an amount")
			}
			amnt, ok := amntI.(*big.Int)
			if !ok {
				return errors.New("Structure is not an amount")
			}
			cV, ok := m.cm.GetContractByNameVersion("validator_service", c.Version)
			if !ok {
				return errors.New("Node contract is not found for version :" + c.Version)
			}
			v, err := m.getValidatorChanged(ctx, m.tr.GetBoundContractCaller(ctx, cV.Addr, cV.Abi), ce.BlockHeight, vID)
			if err != nil {
				return fmt.Errorf("error running validatorChanged  %w", err)
			}
			fee := big.NewInt(0)
			fee.Mul(amnt, v.FeeRate)
			fee.Div(fee, big.NewInt(100))

			var t = int(structs.ValidatorStatisticsTypeUnclaimedRewards)
			params := structs.ValidatorStatisticsParams{
				ValidatorId:      vID.String(),
				StatisticsTypeVS: string(t),
			}
			vldStatistics, err := m.dataStore.GetValidatorStatistics(ctx, params)
			if err != nil {
				return fmt.Errorf("error getting validator statistics %w", err)
			}
			if vldStatistics != nil {
				fee.Add(fee, vldStatistics[0].Amount)
			}
			err = m.updateUnclaimedRewardsValidator(ctx, vID, fee, ce.BlockHeight)
			if err != nil {
				return fmt.Errorf("error updating unclaimed rewards for validator function %w", err)
			}

			bounty := big.NewInt(0)
			bounty.Sub(amnt, fee)
			err = m.updateBountyValidator(ctx, vID, bounty, ce.BlockHeight)
			if err != nil {
				return fmt.Errorf("error updating bounty for validator function %w", err)
			}

		}

	case "delegation_controller":
		/*
			@dev Emitted when a delegation is proposed to a validator.
				event DelegationProposed(
				uint delegationId
			);
			@dev Emitted when a delegation is accepted by a validator.
			event DelegationAccepted(
				uint delegationId
			);
			@dev Emitted when a delegation is cancelled by the delegator.
			event DelegationRequestCanceledByUser(
				uint delegationId
			);
			@dev Emitted when a delegation is requested to undelegate.
			event UndelegationRequested(
				uint delegationId
			);
		*/

		dIDI, ok := ce.Params["delegationId"]
		if !ok {
			return errors.New("structure is not a delegation")
		}
		dID, ok := dIDI.(*big.Int)
		if !ok {
			return errors.New("structure is not a delegation")
		}

		d, err := m.getDelegationChanged(ctx, bc, ce.BlockHeight, dID)
		if err != nil {
			return fmt.Errorf("error running delegationChanged  %w", err)
		}
		d.TransactionHash = ce.TransactionHash
		d.BlockHeight = ce.BlockHeight

		if err := m.dataStore.SaveDelegation(ctx, d); err != nil {
			return fmt.Errorf("error storing delegation %w", err)
		}

		dlgs, _ := m.c.GetHolderDelegations(ctx, bc, ce.BlockHeight, d.Holder)
		accType := "nothing"
		if dlgs != nil && len(dlgs) > 0 {
			accType = "delegator"
		}

		err = m.dataStore.SaveAccount(ctx, structs.Account{
			Address:     d.Holder,
			AccountType: accType,
		})
		if err != nil {
			return fmt.Errorf("error storing account %w", err)
		}

		vs := structs.ValidatorStatisticsParams{
			ValidatorId: d.ValidatorID.String(),
			BlockHeight: ce.BlockHeight,
		}
		if err := m.dataStore.CalculateTotalStake(ctx, vs); err != nil {
			return fmt.Errorf("error storing delegation %w", err)
		}
		ce.BoundType = "delegation"
		ce.BoundID = []big.Int{*dID, *d.ValidatorID}

		if ce.EventName == "DelegationAccepted" || ce.EventName == "UndelegationRequested" {
			dlgs, _ := m.c.GetHolderDelegations(ctx, bc, ce.BlockHeight, d.Holder)
			sumH := big.NewInt(0)
			for _, dlg := range dlgs {
				if dlg.ValidatorID == d.ValidatorID && dlg.State == structs.DelegationStateACCEPTED {
					sumH.Add(sumH, dlg.Amount)
				}
			}

			dlgs, _ = m.c.GetValidatorDelegations(ctx, bc, ce.BlockHeight, d.ValidatorID)
			sumT := big.NewInt(0)
			for _, dlg := range dlgs {
				if dlg.State == structs.DelegationStateACCEPTED {
					sumT.Add(sumT, dlg.Amount)
				}
			}

			var t = int(structs.ValidatorStatisticsTypeBounty)
			params := structs.ValidatorStatisticsParams{
				ValidatorId:      d.ValidatorID.String(),
				StatisticsTypeVS: string(t),
			}
			vldStatistics, err := m.dataStore.GetValidatorStatistics(ctx, params)
			if err != nil {
				return fmt.Errorf("error getting validator statistics %w", err)
			}
			bounty := vldStatistics[0].Amount
			bounty.Mul(bounty, sumH)
			bounty.Div(bounty, sumT)

			err = m.updateUnclaimedRewardsDelegator(ctx, d.Holder, bounty, ce.BlockHeight)
			if err != nil {
				return fmt.Errorf("error updating unclaimed delegator statistics %w", err)
			}
		}
	case "skale_manager":
		/*
			@dev Emitted when bounty is received.
				event BountyReceived(
				uint indexed nodeIndex,
				address owner,
				uint averageDowntime,
				uint averageLatency,
				uint bounty,
				uint previousBlockEvent,
				uint time,
				uint gasSpend
			);
			event BountyGot(
				uint indexed nodeIndex,
				address owner,
				uint averageDowntime,
				uint averageLatency,
				uint bounty,
				uint previousBlockEvent,
				uint time,
				uint gasSpend
			);
		*/
	case "skale_token":
		ce.BoundType = "token"
		if ce.EventName == "transfer" || ce.EventName == "approval" {
			ce, err = standard.DecodeERC20Events(ctx, ce)
			if err != nil {
				return fmt.Errorf("error decoding event ERC20 %w", err)
			}
		} else {
		}

	default:

		m.l.Debug("Unknown event type", zap.String("type", ce.ContractName), zap.Any("event", ce))
	}

	return m.dataStore.SaveContractEvent(ctx, ce)
}

func (m *Manager) getValidatorChanged(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (structs.Validator, error) {

	validator, err := m.c.GetValidator(ctx, bc, blockNumber, validatorID)
	if err != nil {
		return validator, fmt.Errorf("error calling getValidator function %w", err)
	}

	validator.Authorized, err = m.c.IsAuthorizedValidator(ctx, bc, blockNumber, validatorID)
	if err != nil {
		return validator, fmt.Errorf("error calling IsAuthorizedValidator function %w", err)
	}

	return validator, nil
}

func (m *Manager) getDelegationChanged(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (structs.Delegation, error) {

	delegation, err := m.c.GetDelegation(ctx, bc, blockNumber, delegationID)
	if err != nil {
		return delegation, fmt.Errorf("error calling GetDelegation %w", err)
	}

	delegation.State, err = m.c.GetDelegationState(ctx, bc, blockNumber, delegationID)
	if err != nil {
		return delegation, fmt.Errorf("error calling GetDelegationState %w", err)
	}

	return delegation, nil
}

func (m *Manager) updateUnclaimedRewardsValidator(ctx context.Context, validatorID, amount *big.Int, blockHeight uint64, ) error {
	vs := structs.ValidatorStatistics{
		ValidatorId:    validatorID,
		Amount:         amount,
		BlockHeight:    blockHeight,
		StatisticsType: structs.ValidatorStatisticsTypeUnclaimedRewards,
	}
	err := m.dataStore.SaveValidatorStatistics(ctx, vs)
	if err != nil {
		return fmt.Errorf("error saving validator statistics function %w", err)
	}
	err = m.dataStore.UpdateUnclaimedRewards(ctx, vs.ValidatorId, vs.BlockHeight)
	if err != nil {
		return fmt.Errorf("error updating unclaimed rewards function %w", err)
	}
	return nil
}

func (m *Manager) updateClaimedRewardsValidator(ctx context.Context, validatorID, amount *big.Int, blockHeight uint64, ) error {
	vs := structs.ValidatorStatistics{
		ValidatorId:    validatorID,
		Amount:         amount,
		BlockHeight:    blockHeight,
		StatisticsType: structs.ValidatorStatisticsTypeClaimedRewards,
	}
	err := m.dataStore.SaveValidatorStatistics(ctx, vs)
	if err != nil {
		return fmt.Errorf("error saving validator statistics function %w", err)
	}
	err = m.dataStore.UpdateClaimedRewards(ctx, vs.ValidatorId, vs.BlockHeight)
	if err != nil {
		return fmt.Errorf("error updating unclaimed rewards function %w", err)
	}
	return nil
}

func (m *Manager) updateBountyValidator(ctx context.Context, validatorID, amount *big.Int, blockHeight uint64, ) error {
	vs := structs.ValidatorStatistics{
		ValidatorId:    validatorID,
		Amount:         amount,
		BlockHeight:    blockHeight,
		StatisticsType: structs.ValidatorStatisticsTypeBounty,
	}
	err := m.dataStore.SaveValidatorStatistics(ctx, vs)
	if err != nil {
		return fmt.Errorf("error saving validator statistics function %w", err)
	}
	return nil
}

func (m *Manager) updateUnclaimedRewardsDelegator(ctx context.Context, holder common.Address, amount *big.Int, blockHeight uint64, ) error {
	ds := structs.DelegatorStatistics{
		Holder:           holder,
		Amount:           amount,
		BlockHeight:      blockHeight,
		StatisticsTypeDS: structs.DelegatorStatisticsTypeUnclaimedRewards,
	}
	err := m.dataStore.SaveDelegatorStatistics(ctx, ds)
	if err != nil {
		return fmt.Errorf("error saving delegator statistics function %w", err)
	}
	return nil
}

func (m *Manager) updateClaimedRewardsDelegator(ctx context.Context, holder common.Address, amount *big.Int, blockHeight uint64, ) error {
	ds := structs.DelegatorStatistics{
		Holder:           holder,
		Amount:           amount,
		BlockHeight:      blockHeight,
		StatisticsTypeDS: structs.DelegatorStatisticsTypeClaimedRewards,
	}
	err := m.dataStore.SaveDelegatorStatistics(ctx, ds)
	if err != nil {
		return fmt.Errorf("error saving delegator statistics function %w", err)
	}
	return nil
}
