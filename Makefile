.PHONY: install installrace lint init_db test unit_test integration_test \
mod generate

RFLAG =
ifneq ($(RACE),)
RFLAG = -race
export GORACE = "halt_on_error=1"
endif

# Assemble GOBIN until supported: https://github.com/golang/go/issues/23439
CGO =
INSTALLPATH = $(shell go env GOPATH)
ifneq ($(DOCKER),)
CGO = CGO_ENABLED=0
INSTALLPATH = .
endif

ifeq ($(strip $(TEST_PG_URI)),)
TEST_PG_URI = postgres://postgres:postgres@127.0.0.1/atlas_test
endif

# Native Go (CGO_ENABLED=0) is faster for non-esoteric uses of DNS/user, but is
# not currently supported in conjunction with PIE on darwin/amd64:
# https://github.com/golang/go/issues/42459
install:
	for x in $(shell find cmd -mindepth 1 -type d); do $(CGO) go build -o \
	$(INSTALLPATH)/bin/$${x#cmd/} -ldflags="-w" -buildmode=pie ./$${x}; done

	for x in $(shell find tool -mindepth 1 -type d); do $(CGO) go build -o \
	$(INSTALLPATH)/bin/$${x#tool/} -ldflags="-w" -buildmode=pie ./$${x}; done

# Race detector is exclusive of non-cgo and PIE
# https://github.com/golang/go/issues/9918
installrace:
	for x in $(shell find cmd -mindepth 1 -type d); do go build -o \
	$(INSTALLPATH)/bin/$${x#cmd/}.race -ldflags="-w" -race ./$${x}; done

lint:
	cd /tmp && GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck && \
	cd $(CURDIR)
# staticcheck defaults are all,-ST1000,-ST1003,-ST1016,-ST1020,-ST1021,-ST1022
	staticcheck -checks all ./...

	cd /tmp && GO111MODULE=on go get \
	github.com/golangci/golangci-lint/cmd/golangci-lint && cd $(CURDIR)
# unused is included in the newer version of staticcheck above
	golangci-lint run -D staticcheck,unused -E \
	goconst,godot,goerr113,gosec,prealloc,scopelint,unconvert,unparam

init_db:
	cd /tmp && GO111MODULE=on go get -tags postgres \
	github.com/golang-migrate/migrate/v4/cmd/migrate && cd $(CURDIR)
	migrate -path /tmp -database $(TEST_PG_URI)?sslmode=disable drop -f
	migrate -path config/db/atlas -database $(TEST_PG_URI)?sslmode=disable up

test: install installrace lint unit_test integration_test
# -count 1 is the idiomatic way to disable test caching in package list mode
unit_test:
	go test -count=1 -cover $(RFLAG) -cpu 1,4 -tags unit ./...
integration_test: init_db
	go test -count=1 -cover $(RFLAG) -cpu 1,4 -tags integration ./...

mod:
	go get -t -u ./...
	go mod tidy -v
	go mod vendor
# Update atlas.swagger.json at the same time as github.com/thingspect/api
	if [ -f ../api/openapi/atlas.swagger.json ]; then cp -f -v \
	../api/openapi/atlas.swagger.json web/; fi

generate:
	cd /tmp && GO111MODULE=on go get github.com/golang/mock/mockgen && \
	cd $(CURDIR)
	go generate -x ./...
