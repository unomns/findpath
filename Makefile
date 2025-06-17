FILE = map.example.json
ALGO = a-star

PROTO_DIR := ./protos
PROTO_OUT := ./protos/gen

run-cli:
	go run ./cmd/findpath-cli --FILE=$(FILE) --algo=$(ALGO)

cli:
	go build -o bin/findpath-cli ./cmd/findpath-cli

grpc:
	go build -o bin/findpath-grpc ./cmd/findpath-grpc

clean:
	rm ./bin/findpath-*

proto:
	@echo "Generating Protobuf..."
	@protoc \
		-I $(PROTO_DIR)/proto $(PROTO_DIR)/proto/findpath/findpath.proto \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative
