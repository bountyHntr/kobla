proto-compile-pow:
	protoc --go_out=. --go_opt=paths=source_relative ./blockchain/core/pb/blockchain_pow.proto

proto-compile-poa:
	protoc --go_out=. --go_opt=paths=source_relative ./blockchain/core/pb/blockchain_poa.proto
	
build-pow:
	go build --tags pow	-o kobla

build-poa:
	go build --tags poa -o kobla

tests:
	go test --tags pow ./...
	go test --tags pow ./...

config = config.yaml
run:
	./kobla 2> `date +%s`.log --config=$(config)