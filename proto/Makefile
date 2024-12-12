.PHONY: generate clean

# To avoid use of Git submodules, Atlas and Proto must be placed at the same
# filesystem depth
generate:
	cp -Rf ../../proto/protobuf/api ../../proto/protobuf/common protobuf/
	docker compose build --progress=plain --no-cache --pull
	docker compose up --menu=false
	docker compose down
	rm -rf protobuf/api protobuf/common
	cp -f ../../proto/openapi/atlas.swagger.json ../web/

clean:
	find . -name '*.pb*.go' -type f|xargs rm -v
