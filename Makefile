# Makefile

#
# Variables
#
APP = xis-data-aggregator
CMD_MAIN = cmd/$(APP)/main.go

#setup for target system (to mount the container/tmp folder)
LOG_DIR =  ./logs

#
# Tools
#
PROTOC = protoc
PROTOC_GEN_GO = protoc-gen-go
PROTOC_GEN_GO_GRPC = protoc-gen-go-grpc
PROTOC_GEN_DOC = protoc-gen-doc
PROTO_DIR = ./gen/proto

# Docker 
DOCKER_LATEST = $(APP):latest

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
	 	--proto_path=$(PROTO_DIR) \
		--plugin=$(PROTOC_GEN_DOC)=$(shell go env GOPATH)/bin/$(PROTOC_GEN_DOC) \
		--doc_out=./docs \
		--doc_opt=html,index.html \
		$(PROTO_DIR)/*.proto

# Docker targets
.PHONY: docker-build
docker-build:
	docker build -t xis-data-aggregator:latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm -it \
	  --memory="128m" \
	  --name=$(APP) \
	  --volume=$(LOG_DIR):/tmp \
	  -p 8080:8080 \
	  -p 50051:50051 \
	  $(DOCKER_LATEST)
