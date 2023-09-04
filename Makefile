proto-compile-pow:
	protoc --go_out=. --go_opt=paths=source_relative ./blockchain/core/pb/blockchain_pow.proto

proto-compile-poa:
	protoc --go_out=. --go_opt=paths=source_relative ./blockchain/core/pb/blockchain_poa.proto
	
build-pow:
	go build --tags pow	-o kobla

build-poa:
	go build --tags poa -o kobla

config = config.yaml
run:
	./kobla 2> `date +%s`.log --config=$(config)