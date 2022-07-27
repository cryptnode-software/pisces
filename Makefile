CURRENT_DIRECTORY = $(shell pwd)

.PHONY: show-help
help: show-help

.PHONY: generate-js
generate-js:
		protoc -I $(shell pwd)/proto/ \
    	--plugin=protoc-gen-grpc-web=./client/node_modules/.bin/protoc-gen-grpc-web \
    	--js_out=import_style=commonjs,binary:$(shell pwd)/client/src/proto \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:$(shell pwd)/client/src/proto \
    	$(shell pwd)/proto/main.proto