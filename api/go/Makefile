.PHONY: generate deps message pwt

generate: deps message pwt

deps:
# Force GO111MODULE on in the absence of a go.mod file
	GO111MODULE=on go get google.golang.org/protobuf/cmd/protoc-gen-go

message:
	protoc -I=../protobuf --go_out=. --go_opt=paths=source_relative \
	message/validator_in.proto \
	message/validator_out.proto

pwt:
	protoc -I=../protobuf --go_out=. --go_opt=paths=source_relative \
	pwt/claim.proto
