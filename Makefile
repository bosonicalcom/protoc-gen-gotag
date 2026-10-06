LOCAL_PATH = $(shell pwd)
PROTOC_INCLUDE ?= /usr/local/include
EDITIONS_PROTOS = example/editions2023/editions2023.proto example/editions2024/editions2024.proto

.PHONY: example proto install gen-tag test testdata

example: proto install
	protoc -I ${PROTOC_INCLUDE} \
	-I ${LOCAL_PATH} \
	--gotag_out=paths=source_relative,xxx="graphql+\"-\" bson+\"-\"":. example/example.proto
	protoc -I ${PROTOC_INCLUDE} \
	-I ${LOCAL_PATH} \
	--gotag_out=paths=source_relative:. ${EDITIONS_PROTOS}

proto:
	protoc -I ${PROTOC_INCLUDE} \
	-I ${LOCAL_PATH} \
	--go_out=paths=source_relative:. example/example.proto
	protoc -I ${PROTOC_INCLUDE} \
	-I ${LOCAL_PATH} \
	--go_out=paths=source_relative:. ${EDITIONS_PROTOS}

# Untagged protoc-gen-go output and descriptors used by module/editions_test.go
testdata:
	protoc -I ${PROTOC_INCLUDE} \
	-I ${LOCAL_PATH} \
	--go_out=paths=source_relative:module/testdata ${EDITIONS_PROTOS}
	protoc -I ${PROTOC_INCLUDE} \
	-I ${LOCAL_PATH} \
	--include_imports --include_source_info \
	--descriptor_set_out=module/testdata/editions.binpb ${EDITIONS_PROTOS}

install:
	go install .

gen-tag:
	buf generate
	buf generate --template=buf.gen.tag.yaml
	buf generate --template=buf.gen.debug.yaml --path tagger

test:
	go test ./...
