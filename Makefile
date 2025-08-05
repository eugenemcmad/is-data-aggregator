# Makefile

#
# Variables
#
APP = xis-data-aggregator
BIN = $(APP)
BR = `git name-rev --name-only HEAD`
VER = `git describe --tags --abbrev=0`
COMMIT = `git rev-parse --short HEAD`
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
FORMAT = v$(VER)-$(COMMIT)-$(TIMESTM)
DOCTAG = $(VER)-$(BR)
PROTO_DIR = gen/proto
CMD_MAIN = cmd/$(APP)/main.go

#
# Tools
#
PROTOC = protoc
PROTOC_GEN_GO = protoc-gen-go
PROTOC_GEN_GO_GRPC = protoc-gen-go-grpc

#
# Targets
#

swag: ## Generates swagger files
	@echo "Generating swagger files..."
	@init swag init -g $(CMD_MAIN)

proto: ## Generates Go code from .proto files
	@echo "Generating Protobuf Go code..."
	$(PROTOC) \
	  --proto_path=gen \
	  --go_out=. \
	  --go-grpc_out=. \
	  $(PROTO_DIR)/*.proto

grpcdoc: ## Generates swagger files from .proto files using protoc-gen-doc and
	$(PROTOC) \
	 	--proto_path=./gen/proto \
		#--plugin=protoc-gen-doc.exe=$(shell go env GOPATH)/bin/protoc-gen-doc \
		--plugin=protoc-gen-doc=$(shell go env GOPATH)/bin/protoc-gen-doc \
		--doc_out=./docs \
		--doc_opt=html,index.html \
		./gen/proto/*.proto


unit-test:
	CGO_ENABLED=0 go test ./...


build:
	CGO_ENABLED=0 go build -o $(BIN) -ldflags "-X main.BuildVersion=$(FORMAT)"

build-image:
	sudo docker build -t $(APP):$(DOCTAG) .
