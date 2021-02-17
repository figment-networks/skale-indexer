package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
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
	nvIndex    map[NV]ContractsContents

	globalEvents map[string]abi.Event
}
type NV struct {
	Name    string
	Version string
}

func NewManager() *Manager {
	return &Manager{
		contractsTable: make(map[common.Address]ContractsContents),
		nvIndex:        make(map[NV]ContractsContents),
		globalEvents:   make(map[string]abi.Event),
	}
}

type Contracts struct {
	ccs map[common.Address][]ContractsContents
}

func NewContracts() *Contracts {
	return &Contracts{ccs: make(map[common.Address][]ContractsContents)}
}

func (c *Contracts) GetByAddress(a common.Address) (cc ContractsContents, ok bool) {
	s, ok := c.ccs[a]
	if !ok {
		return cc, false
	}

	return s[0], true
}

func (c *Contracts) GetAddresses() (addrs []common.Address) {
	for _, a := range c.ccs {
		addrs = append(addrs, a[0].Addr)
	}
	return addrs
}

func (c *Contracts) GetAllVersions(a common.Address) (cc []ContractsContents, ok bool) {
	s, ok := c.ccs[a]
	if !ok {
		return cc, false
	}

	return s, true
}

func (c *Contracts) SetAllVersions(a common.Address, cc []ContractsContents) {
	c.ccs[a] = cc
}

type ContractsContents struct {
	Name    string
	Addr    common.Address
	Abi     abi.ABI
	Bound   *bind.BoundContract
	Version string
}

func (m *Manager) AddGlobalEvents(readr io.Reader) error {
	parsedABI, err := abi.JSON(readr)
	if err != nil {
		return err
	}

	if parsedABI.Events != nil {
		m.lookupLock.Lock()
		for k, ev := range parsedABI.Events {
			m.globalEvents[k] = ev
		}
		defer m.lookupLock.Unlock()
	}

	return nil
}

func (m *Manager) GetContractsByNames(names []string) *Contracts {

	ccs := NewContracts()
	for _, n := range names {
		contr := []ContractsContents{}
		for nv, c := range m.nvIndex {
			if nv.Name == n {
				contr = append(contr, c)
			}
		}
		if len(contr) > 0 {
			sort.Slice(contr, func(i, j int) bool {
				vi := strings.Split(contr[i].Version, ".")
				vj := strings.Split(contr[j].Version, ".")
				v1i, err := strconv.Atoi(vi[0])
				if err != nil {
					return false
				}
				v1j, err := strconv.Atoi(vj[0])
				if err != nil {
					return false
				}
				if v1i != v1j {
					return v1i > v1j
				}

				v2i, err := strconv.Atoi(vi[1])
				if err != nil {
					return false
				}
				v2j, err := strconv.Atoi(vj[1])
				if err != nil {
					return false
				}
				if v2i != v2j {
					return v2i > v2j
				}

				v3i, err := strconv.Atoi(vi[2])
				if err != nil {
					return false
				}
				v3j, err := strconv.Atoi(vj[2])
				if err != nil {
					return false
				}

				return v3i > v3j
			})
			ccs.SetAllVersions(contr[0].Addr, contr)
		}

	}

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
					err = m.getContracts(abiF, d.Name())
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

func (m *Manager) attachGlobalEvents(inputABI map[string]abi.Event) {
	for k, ev := range m.globalEvents {
		inputABI[k] = ev
	}
}

func (m *Manager) LoadContract(name, addr, version string, abiContents abi.ABI) error {
	cc := ContractsContents{Name: name, Abi: abiContents, Version: version}
	cc.Addr = common.HexToAddress(addr)

	m.attachGlobalEvents(abiContents.Events)

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
