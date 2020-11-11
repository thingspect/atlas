export GORACE = "halt_on_error=1"
RFLAG =
ifneq ($(RACE),)
RFLAG = -race
endif

lint:
	cd /tmp && GO111MODULE=on go get honnef.co/go/tools/cmd/staticcheck && \
	cd $(CURDIR)
# staticcheck defaults are all,-ST1000,-ST1003,-ST1016,-ST1020,-ST1021,-ST1022
#	staticcheck -checks all -unused.whole-program ./...
	staticcheck -checks all ./...
	cd /tmp && GO111MODULE=on go get \
	github.com/golangci/golangci-lint/cmd/golangci-lint && cd $(CURDIR)
# unused is included in the newer version of staticcheck above
	golangci-lint run -D staticcheck,unused -E \
	goconst,godot,goerr113,gosec,prealloc,unconvert,unparam

# -count 1 is the idiomatic way to disable test caching in package list mode
test: lint unit_test integration_test
unit_test:
	go test -count=1 -cover -race -cpu 1,4 -tags unit ./...
integration_test:
	go test -count=1 -cover $(RFLAG) -cpu 1,4 -tags integration ./...
