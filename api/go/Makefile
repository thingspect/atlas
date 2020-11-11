.PHONY: generate deps message mod

generate: deps message mod

deps:
	GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go

message:
	protoc -I=../protobuf --go_out=. --go_opt=paths=source_relative \
	message/message.proto

mod:
	go get -u ./... && go mod tidy
