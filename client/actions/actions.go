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
	StoreValidator(ctx context.Context, height uint64, v structures.Validator) error
	StoreEvent(ctx context.Context, v structures.ContractEvent) error
}

type Manager struct {
	s Store
	c Call
}

func (m *Manager) StoreEvent(ctx context.Context, ev structures.ContractEvent) error {

	// some more magic in will be here in future

	return m.s.StoreEvent(ctx, ev)
}

func (m *Manager) ValidatorChanged(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) error {

	validator, err := m.c.GetValidator(ctx, bc, blockNumber, validatorID)
	if err != nil {
		return fmt.Errorf("error calling getValidator function %w", err)
	}

	validator.Authorized, err = m.c.IsAuthorizedValidator(ctx, bc, blockNumber, validatorID)
	if err != nil {
		return fmt.Errorf("error calling IsAuthorizedValidator function %w", err)
	}

	return nil
}

func (m *Manager) AfterEventLog(ctx context.Context, bc *bind.BoundContract, ce structures.ContractEvent) error {
	// c contract.ContractsContents,
	// client bind.ContractCaller,

	//	blockNumber uint64,
	//	mapped map[string]interface{}) error {
	// TODO(lukanus): Cache is somehow
	switch ce.ContractName {
	case "validator_service":
		if ce.Type == "NodeAddressWasAdded" || ce.Type == "NodeAddressWasRemoved" {
			break
		}

		vID, ok := ce.Params["validatorId"]
		if !ok {
			//		logger.Error("Calling getValidator %+v", zap.Any("validatorID", vID))
			return errors.New("Structure is not a validator")
		}

		v, err := m.c.GetValidator(ctx, bc, ce.Height, vID.(*big.Int))
		if err != nil {
			return errors.New("Structure is not a validator")
		}
		m.s.StoreValidator(ctx, ce.Height, v)

	case "nodes":
		/*	event Slash(
				uint validatorId,
				uint amount
			);
			event Forgive(
				address wallet,
				uint amount
			);*/
	case "punisher":
		/*	event Slash(
				uint validatorId,
				uint amount
			);
			event Forgive(
				address wallet,
				uint amount
			);*/
	case "distributor":
		/**
		 * @dev Emitted when bounty is withdrawn.
		event WithdrawBounty(
			address holder,
			uint validatorId,
			address destination,
			uint amount
		);
		* @dev Emitted when a validator fee is withdrawn.
		event WithdrawFee(
			uint validatorId,
			address destination,
			uint amount
		);
		 * @dev Emitted when bounty is distributed.
		event BountyWasPaid(
			uint validatorId,
			uint amount
		);
		*/
		vID, ok := ce.Params["validatorId"]
		if !ok {
			return errors.New("Structure is not a validator")
		}
		earned, endMonth, err := m.c.GetEarnedFeeAmountOf(ctx, bc, ce.Height, vID.(*big.Int))
		if err != nil {
			return fmt.Errorf("error calling getEarnedFeeAmountOf function %w", err)
		}
	//	logger.Debug("got distributor", zap.Any("earned", earned), zap.Any("endMonth", endMonth))
	case "delegation_controller":
		dID, ok := ce.Params["delegationId"]
		if !ok {
			return errors.New("Structure is not a delegation")
		}

		d, err := m.c.GetDelegation(ctx, bc, ce.Height, dID.(*big.Int))
		if err != nil {
			return fmt.Errorf("error calling getDelegation function %w", err)
		}
	//	logger.Debug("got delegation", zap.Any("delegation", d))
	case "skale_manager":

		/**
		     * @dev Emitted when bounty is received.

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
