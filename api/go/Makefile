.PHONY: generate deps message

generate: deps message

deps:
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go

message:
	protoc -I=../protobuf --go_out=. --go_opt=paths=source_relative \
	message/validator_in.proto message/validator_out.proto
