
LDFLAGS      := -w -s
MODULE       := github.com/figment-networks/skale-indexer
VERSION_FILE ?= ./VERSION


# Git Status
GIT_SHA ?= $(shell git rev-parse --short HEAD)

ifneq (,$(wildcard $(VERSION_FILE)))
VERSION ?= $(shell head -n 1 $(VERSION_FILE))
else
VERSION ?= n/a
endif

all: generate build

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build: LDFLAGS += -X $(MODULE)/cmd/skale-indexer/config.Timestamp=$(shell date +%s)
build: LDFLAGS += -X $(MODULE)/cmd/skale-indexer/config.Version=$(VERSION)
build: LDFLAGS += -X $(MODULE)/cmd/skale-indexer/config.GitSHA=$(GIT_SHA)
build:
	$(info building indexer binary as ./indexer)
	go build -o indexer -ldflags '$(LDFLAGS)' ./cmd/skale-indexer


.PHONY: build-migration
build-migration:
	$(info building migration binary as ./migration)
	go build -o migration ./cmd/skale-indexer-migration

.PHONY: pack-release
pack-release:
	$(info preparing release)
	@mkdir -p ./release
	@make build-migration
	@mv ./migration ./release/migration
	@make build
	@cp -R ./cmd/skale-indexer-migration/migrations ./release/
	@zip -r indexer ./release
	@rm -rf ./release

.PHONY: generate-types
generate-types:
	@mkdir -p ./install
	if [ ! -d ./install/skale-network ]; then \
		git clone https://github.com/skalenetwork/skale-network ./install/skale-network; \
	else \
		cd ./install/skale-network && git pull; \
	fi

	abigen --combined-json ./install/skale-network/releases/mainnet/skale-manager/1.5.0/abi.json --out a.go --pkg api --lang go

	cd ./install/skale-manager/contracts/delegation/ && solc  --allow-paths .. --abi ./DelegationController.sol

.PHONY: install-deps
install-deps:
	mkdir -p ./install
	git clone https://github.com/ethereum/go-ethereum.git ./install/go-ethereum
	$(MAKE) -C ./install/go-ethereum
	make all

