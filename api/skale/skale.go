package skale

import "golang.org/x/time/rate"

type EthereumNodeType uint8

const (
	ENTArchive EthereumNodeType = iota
	ENTRecent
)

// Caller caller for SKALE api functions
// The only reason this exists is to be used as an interface
// binding all it's methods
type Caller struct {
	NodeType                  EthereumNodeType
	validatorDelegationsCache map[uint64]ValidatorDelegationsCache // validatorID
	rateLimiter               *rate.Limiter
}

type ValidatorDelegationsCache struct {
	LastID      uint64
	Length      uint64
	Delegations []uint64
}

func NewCaller(NodeType EthereumNodeType, requestsPerSecond float64) *Caller {

	rateLimiter := rate.NewLimiter(rate.Limit(requestsPerSecond), int(requestsPerSecond))
	return &Caller{
		NodeType:                  NodeType,
		rateLimiter:               rateLimiter,
		validatorDelegationsCache: make(map[uint64]ValidatorDelegationsCache),
	}
}
