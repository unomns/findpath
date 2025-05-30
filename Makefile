testFile = map.json
algorithm = a

run:
	go run cmd/main.go --file=$(testFile) --algo=$(algorithm)