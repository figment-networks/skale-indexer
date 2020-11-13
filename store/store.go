package store

import (
	"../structs"
	"../types"
	"context"
)

type DBDriver interface {
	DelegationStore
}

type DataStore interface {
	DelegationStore
}

type DelegationStore interface {
	SaveOrUpdateDelegation(ctx context.Context, delegation structs.Delegation) error
	SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error
	GetDelegationById(ctx context.Context, id types.ID) (res structs.Delegation, err error)
	GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error)
	GetDelegationsByValidatorId(ctx context.Context, validatorId uint64) (delegations []structs.Delegation, err error)
}

type Store struct {
	driver DBDriver
}

func New(driver DBDriver) *Store {
	return &Store{driver: driver}
}

func (s *Store) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (s *Store) SaveOrUpdateDelegation(ctx context.Context, delegation structs.Delegation) error {
	return s.driver.SaveOrUpdateDelegation(ctx, delegation)
}

func (s *Store) SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error {
	return s.driver.SaveOrUpdateDelegations(ctx, delegations)
}

func (s *Store) GetDelegationById(ctx context.Context, id types.ID) (res structs.Delegation, err error) {
	return s.driver.GetDelegationById(ctx, id)
}

func (s *Store) GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error) {
	return s.driver.GetDelegationsByHolder(ctx, holder)
}

func (s *Store) GetDelegationsByValidatorId(ctx context.Context, validatorId uint64) (delegations []structs.Delegation, err error) {
	return s.driver.GetDelegationsByValidatorId(ctx, validatorId)
}
