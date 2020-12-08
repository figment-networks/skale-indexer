package actions

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/client/structures"
)

var implementedEvents = []string{"delegation_controller", "validator_service", "nodes", "distributor", "punisher", "skale_manager", "bounty", "bounty_v2"}

type Call interface {
	// Validator
	IsAuthorizedValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (isAuthorized bool, err error)
	GetValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (v structures.Validator, err error)

	// Nodes
	GetValidatorNodes(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (nodes []structures.Node, err error)
	GetNode(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (n structures.Node, err error)
	GetNodeNextRewardDate(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (t time.Time, err error)

	// Distributor
	GetEarnedFeeAmountOf(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (earned, endMonth *big.Int, err error)

	// Delegation
	GetPendingDelegationsTokens(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, holderAddress common.Address) (amount *big.Int, err error)
	GetDelegation(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (d structures.Delegation, err error)
	GetDelegationState(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (ds structures.DelegationState, err error)
	GetValidatorDelegations(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (delegations []structures.Delegation, err error)
}

type Store interface {
	StoreEvent(ctx context.Context, v structures.ContractEvent) error

	StoreValidator(ctx context.Context, height uint64, t time.Time, v structures.Validator) error
	StoreDelegation(ctx context.Context, height uint64, t time.Time, d structures.Delegation) error

	StoreNode(ctx context.Context, height uint64, t time.Time, v structures.Node) error
	StoreValidatorNodes(ctx context.Context, height uint64, t time.Time, nodes []structures.Node) error
}

type Calculator interface {
	ValidatorParams(ctx context.Context, height uint64, vID *big.Int) error
	DelegationParams(ctx context.Context, height uint64, dID *big.Int) error
}

type Manager struct {
	s    Store
	c    Call
	calc Calculator
}

func NewManager(c Call, s Store, calc Calculator) *Manager {
	return &Manager{c: c, s: s, calc: calc}
}

func (m *Manager) StoreEvent(ctx context.Context, ev structures.ContractEvent) error {

	// some more magic in will be here in future
	return m.s.StoreEvent(ctx, ev)
}

func (m *Manager) GetImplementedEventsNames() []string {
	return implementedEvents
}

func (m *Manager) validatorChanged(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (structures.Validator, error) {

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

func (m *Manager) AfterEventLog(ctx context.Context, bc *bind.BoundContract, ce structures.ContractEvent) error {
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

		vID, ok := ce.Params["validatorId"]
		if !ok {
			return errors.New("Structure is not a validator")
		}

		v, err := m.validatorChanged(ctx, bc, ce.Height, vID.(*big.Int))
		if err != nil {
			return errors.New("Structure is not a validator")
		}

		if err = m.s.StoreValidator(ctx, ce.Height, ce.Time, v); err != nil {
			return fmt.Errorf("error storing delegation %w", err)
		}

		if ce.Type == "NodeAddressWasAdded" || ce.Type == "NodeAddressWasRemoved" {

			nodes, err := m.c.GetValidatorNodes(ctx, bc, ce.Height, vID.(*big.Int))
			if err != nil {
				return errors.New("Error getting validator nodes")
			}

			if err := m.s.StoreValidatorNodes(ctx, ce.Height, ce.Time, nodes); err != nil {
				return fmt.Errorf("error storing validator nodes %w", err)
			}
		}

		if err = m.calc.ValidatorParams(ctx, ce.Height, vID.(*big.Int)); err != nil {
			return fmt.Errorf("error calculating validator params %w", err)
		}
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

		vID, ok := ce.Params["nodeIndex"]
		if !ok {
			return errors.New("Structure is not a validator")
		}
		n, err := m.c.GetNode(ctx, bc, ce.Height, vID.(*big.Int))
		if err != nil {
			return errors.New("Structure is not a validator")
		}

		if err = m.s.StoreNode(ctx, ce.Height, ce.Time, n); err != nil {
			return fmt.Errorf("error storing delegation %w", err)
		}
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
		/*
			vID, ok := ce.Params["validatorId"]
			if !ok {
				return errors.New("Structure is not a validator")
			}
			earned, endMonth, err := m.c.GetEarnedFeeAmountOf(ctx, bc, ce.Height, vID.(*big.Int))
			if err != nil {
				return fmt.Errorf("error calling getEarnedFeeAmountOf function %w", err)
			}
		*/
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

		dID, ok := ce.Params["delegationId"]
		if !ok {
			return errors.New("Structure is not a delegation")
		}

		d, err := m.c.GetDelegation(ctx, bc, ce.Height, dID.(*big.Int))
		if err != nil {
			return fmt.Errorf("error calling getDelegation function %w", err)
		}

		if err := m.s.StoreDelegation(ctx, ce.Height, ce.Time, d); err != nil {
			return fmt.Errorf("error storing delegation %w", err)
		}

		if err := m.calc.DelegationParams(ctx, ce.Height, dID.(*big.Int)); err != nil {
			return fmt.Errorf("error calculating delegation params %w", err)
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

	}
	return nil

}
