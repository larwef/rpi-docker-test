include config.env

BUILD_VERSION=$(shell git describe --tag --always | sed -e 's,^v,,')
TARGET=target
APP=my-test-app

# ----------------------------------- Proto -----------------------------------
PROTODIR := ./
PROTOC := protoc
INCLUDES := -I $(PROTODIR)
PROTO_OUT = ./
PROTO_SRC = $(shell find . -name '*.proto')

.PHONY: proto
proto: $(patsubst %.proto, %.pb.go, $(PROTO_SRC))
%.pb.go: %.proto
	$(PROTOC) $(INCLUDES) \
	--go_out $(PROTO_OUT) --go_opt paths=source_relative \
	--go-grpc_out $(PROTO_OUT) --go-grpc_opt paths=source_relative \
	$<

# ------------------------------------- Go -------------------------------------
.PHONY: build
build: build-app docker-build

build-app:
	GOOS=linux go build -ldflags "-X main.version=$(BUILD_VERSION)" -o $(TARGET)/$(APP).bin cmd/$(APP)/main.go

.PHONY: test
test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

# ---------------------------------- Docker -----------------------------------
.PHONY: docker
docker: docker-build docker-push

docker-build:
	docker build -t $(REPOSITORY)/$(APP):$(BUILD_VERSION) --build-arg target=$(TARGET)/$(APP).bin -f build/package/Dockerfile .

docker-push:
	docker push $(REPOSITORY)/$(APP):$(BUILD_VERSION)

compose-up:
	IMAGE=$(REPOSITORY)/$(APP):$(BUILD_VERSION) docker-compose -f deployments/docker-compose/docker-compose.yml up

compose-down:
	IMAGE=$(REPOSITORY)/$(APP):$(BUILD_VERSION) docker-compose -f deployments/docker-compose/docker-compose.yml down

# ------------------------------------ SQL -------------------------------------
sqlc:
	sqlc generate 