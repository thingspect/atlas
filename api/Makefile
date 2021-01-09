.PHONY: generate clean

# To avoid use of Git submodules, Atlas and API must be placed at the same
# filesystem depth
generate:
	cp -Rf ../../api/protobuf/common protobuf/
	docker-compose build --no-cache --pull
	docker-compose up --remove-orphans
	rm -rf protobuf/common
	cp -f ../../api/openapi/atlas.swagger.json ../web/

clean:
	find . -name '*.pb*.go' -type f|xargs rm -v
