package types

import (
	"crypto/sha256"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

const (
	HashBits  = 256
	HashBytes = 32
)

var EmptyHash = Hash{}

type Hash [HashBytes]byte

func NewHash(data []byte) Hash {
	hash := sha256.Sum256(data)
	return hash
}

func HashFromSlice(hashData []byte) (hash Hash) {
	if len(hashData) != HashBytes {
		log.Panicf("invalid hash length: %d; %s", len(hashData), string(hashData))
	}

	copy(hash[:], hashData)
	return
}

func HashFromHex(hexData string) (hash Hash) {
	data, err := hex.DecodeString(hexData)
	if err != nil {
		log.Panicf("invalid hex data: %s", hexData)
	}

	return HashFromSlice(data)
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}
