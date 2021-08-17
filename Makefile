.PHONY: install lint init_db test unit_test integration_test mod generate

# Non-cgo DNS is more reliable and faster for non-esoteric uses of resolv.conf
export CGO_ENABLED = 0
RFLAG = -buildmode=pie

# Race detector is exclusive of non-cgo and PIE
# https://github.com/golang/go/issues/6508
ifneq ($(RACE),)
export CGO_ENABLED = 1
RFLAG = -race
export GORACE = halt_on_error=1
endif

# Assemble GOBIN until supported: https://github.com/golang/go/issues/23439
INSTALLPATH = $(shell go env GOPATH)
ifneq ($(DOCKER),)
INSTALLPATH = .
endif

ifeq ($(strip $(TEST_REDIS_HOST)),)
TEST_REDIS_HOST = 127.0.0.1
endif

ifeq ($(strip $(TEST_PG_URI)),)
TEST_PG_URI = postgres://postgres:postgres@127.0.0.1/atlas_test
endif

install:
	for x in $(shell find cmd -mindepth 1 -type d); do go build -ldflags="-w" \
	-o $(INSTALLPATH)/bin/$${x#cmd/} $(RFLAG) ./$${x}; done

	for x in $(shell find tool -mindepth 1 -type d); do go build -ldflags="-w" \
	-o $(INSTALLPATH)/bin/$${x#tool/} $(RFLAG) ./$${x}; done

lint:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -version
# staticcheck defaults are all,-ST1000,-ST1003,-ST1016,-ST1020,-ST1021,-ST1022
	staticcheck -checks all ./...

	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint --version
	golangci-lint run -D staticcheck -E durationcheck,errname,exportloopref \
	-E forcetypeassert,goconst,godot,goerr113,gofumpt,gosec,nlreturn,prealloc \
	-E unconvert,unparam --exclude-use-default=false

init_db:
	echo FLUSHALL|nc -w 2 $(TEST_REDIS_HOST) 6379

	go install -tags postgres \
	github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -path /tmp -database $(TEST_PG_URI)?sslmode=disable drop -f
	migrate -path config/db/atlas -database $(TEST_PG_URI)?sslmode=disable up

test: install lint unit_test integration_test
# -count 1 is the idiomatic way to disable test caching in package list mode
unit_test:
	go test -count=1 -cover -cpu 1,4 -failfast $(RFLAG) -tags unit ./...
integration_test: init_db
	go test -count=1 -cover -cpu 1,4 -failfast $(RFLAG) -tags integration ./...

mod:
	go get -t -u ./... || true
	go mod tidy -v
	go mod vendor
# Update atlas.swagger.json at the same time as github.com/thingspect/api
	if [ -f ../api/openapi/atlas.swagger.json ]; then cp -f -v \
	../api/openapi/atlas.swagger.json web/; fi

generate:
	go install github.com/golang/mock/mockgen@latest
	mockgen -version
	go generate -x ./...
