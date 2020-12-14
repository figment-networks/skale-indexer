package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Manager struct {
	lookupLock     sync.RWMutex
	contractsTable map[common.Address]ContractsContents

	filterLock sync.RWMutex
	// subset of contracts by names fnv64 hash
	addressFilter map[uint64]map[common.Address]ContractsContents
	nvIndex       map[NV]ContractsContents
}
type NV struct {
	Name    string
	Version string
}

func NewManager() *Manager {
	return &Manager{
		contractsTable: make(map[common.Address]ContractsContents),
		addressFilter:  make(map[uint64]map[common.Address]ContractsContents),
		nvIndex:        make(map[NV]ContractsContents),
	}
}

type ContractsContents struct {
	Name    string
	Addr    common.Address
	Abi     abi.ABI
	Bound   *bind.BoundContract
	Version string
}

func (m *Manager) GetContractsByContractNames(names []string) (ccs map[common.Address]ContractsContents) {
	hash := fnv.New64a()
	for _, n := range names {
		hash.Write([]byte(n))
	}
	hSum := hash.Sum64()

	m.filterLock.RLock()
	addr, ok := m.addressFilter[hSum]
	m.filterLock.RUnlock()
	if ok {
		return addr
	}

	ccs = make(map[common.Address]ContractsContents)
	for _, n := range names {
		for _, c := range m.contractsTable {
			if c.Name == n {
				cs := c
				ccs[c.Addr] = cs
			}
		}
	}

	m.filterLock.Lock()
	m.addressFilter[hSum] = ccs
	m.filterLock.Unlock()

	return ccs
}

// LoadContractsFromDir loads abi contracts specifically from skale-network repo path
func (m *Manager) LoadContractsFromDir(inputFolder string) error {
	directories, err := ioutil.ReadDir(inputFolder)
	if err != nil {
		return fmt.Errorf("error reading directory (%s) %w ", inputFolder, err)
	}

	for _, d := range directories {
		if d.IsDir() {
			internalDir := path.Join(inputFolder, d.Name())
			internalFiles, err := ioutil.ReadDir(internalDir)
			if err != nil {
				return fmt.Errorf("error reading directory internal (%s) %w ", internalDir, err)
			}
			for _, f := range internalFiles {
				if strings.Contains(f.Name(), "abi") && !strings.Contains(f.Name(), "token") {
					filePath := path.Join(inputFolder, d.Name(), f.Name())
					abiF, err := os.Open(filePath)
					if err != nil {
						abiF.Close()
						return fmt.Errorf("error reading file (%s) %w ", filePath, err)
					}

					err = m.getContracts(abiF, inputFolder)
					if err != nil {
						abiF.Close()

						return fmt.Errorf("error creating  (%s) %w ", internalDir, err)
					}
					abiF.Close()
				}
			}
		}
	}
	return nil
}

func (m *Manager) getContracts(abiF io.Reader, version string) error {

	abiP := make(map[string]AbiPair)
	dec := json.NewDecoder(abiF)

	if err := dec.Decode(&abiP); err != nil {
		return err
	}

	// (lukanus): crosslink data
	for name, v := range abiP {
		if strings.HasSuffix(name, "_abi") {
			realName := name[0 : len(name)-4]
			addr := abiP[realName+"_address"]

			m.LoadContract(realName, addr.Address, version, v.ABI)
		}
	}
	return nil
}

func (m *Manager) LoadContract(name, addr, version string, abiContents abi.ABI) error {
	cc := ContractsContents{Name: name, Abi: abiContents, Version: version}
	cc.Addr = common.HexToAddress(addr)

	m.lookupLock.Lock()
	m.contractsTable[cc.Addr] = cc
	m.nvIndex[NV{name, version}] = cc
	m.lookupLock.Unlock()
	return nil
}

func (m *Manager) GetContract(addr common.Address) (ContractsContents, bool) {
	m.lookupLock.RLock()
	defer m.lookupLock.RUnlock()
	cc, ok := m.contractsTable[addr]
	return cc, ok
}

func (m *Manager) GetContractByNameVersion(name, version string) (ContractsContents, bool) {
	m.lookupLock.RLock()
	defer m.lookupLock.RUnlock()
	cc, ok := m.nvIndex[NV{name, version}]
	return cc, ok
}

func ConvertToAddress(b []byte) (a common.Address) {
	if len(b) > len(a) {
		b = b[len(b)-common.AddressLength:]
	}
	copy(a[common.AddressLength-len(b):], b)
	return a
}

type AbiPair struct {
	ABI     abi.ABI
	Address string
}

func (ap *AbiPair) UnmarshalJSON(b []byte) error {

	if err := json.Unmarshal(b, &ap.ABI); err != nil {
		ap.Address = string(
			bytes.Replace(b, []byte(`"`), []byte(""), -1),
		)
	}
	return nil
}
