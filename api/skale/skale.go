package skale

type EthereumNodeType uint8

const (
	ENTArchive EthereumNodeType = iota
	ENTRecent
)

// Caller caller for SKALE api functions
// The only reason this exists is to be used as an interface
// binding all it's methods
type Caller struct {
	NodeType EthereumNodeType
}
