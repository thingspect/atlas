.PHONY: generate deps message token

generate: deps message token

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

message:
	protoc -I=../protobuf -I=/tmp/googleapis -I=/tmp/pgv \
	--go_out=. --go_opt=paths=source_relative \
	message/decoder_in.proto \
	message/validator_in.proto \
	message/validator_out.proto

token:
	protoc -I=../protobuf --go_out=. --go_opt=paths=source_relative \
	token/web.proto \
	token/page.proto
