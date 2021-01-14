# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.14 AS buildIndexer

WORKDIR /go/src/github.com/figment-networks/skale-indexer/

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY .git .git
COPY ./Makefile ./Makefile
COPY ./api ./api
COPY ./client ./client
COPY ./scraper ./scraper
COPY ./store ./store
COPY ./cmd/skale-indexer ./cmd/skale-indexer

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux

RUN \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$(git rev-parse HEAD) \
  make build

# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.10 AS releaseIndexer

WORKDIR /app/indexer
COPY --from=buildIndexer /go/src/github.com/figment-networks/skale-indexer/indexer /app/indexer/indexer

RUN chmod a+x ./indexer
CMD ["./indexer"]
