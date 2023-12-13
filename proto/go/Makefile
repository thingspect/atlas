.PHONY: generate deps message token

generate: deps message token

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

message:
	protoc -I=../protobuf -I=/tmp/googleapis -I=/tmp/grpc-gateway -I=/tmp/pgv \
	--go_out=. --go_opt=paths=source_relative \
	message/thingspect_decoder_in.proto \
	message/thingspect_eventer_out.proto \
	message/thingspect_validator_in.proto \
	message/thingspect_validator_out.proto

token:
	protoc -I=../protobuf --go_out=. --go_opt=paths=source_relative \
	token/thingspect_web.proto \
	token/thingspect_page.proto
