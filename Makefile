testFile = map.json
algorithm = a

run:
	go run cmd/main.go --file=$(testFile) --algo=$(algorithm)

generate:
	protoc -I protos/proto protos/proto/findpath/findpath.proto \
	  --go_out=./protos/gen --go_opt=paths=source_relative \
	  --go-grpc_out=./protos/gen --go-grpc_opt=paths=source_relative
