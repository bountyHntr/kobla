proto_compile_pow:
	protoc --go_out=. --go_opt=paths=source_relative ./blockchain/core/pb/blockchain_pow.proto

proto_compile_poa:
	protoc --go_out=. --go_opt=paths=source_relative ./blockchain/core/pb/blockchain_poa.proto