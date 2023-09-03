package types

import (
	"github.com/btcsuite/btcutil/base58"
	streebog "github.com/ddulesov/gogost/gost34112012512"

	log "github.com/sirupsen/logrus"
)

const (
	HashBits  = 512
	HashBytes = HashBits / 8
)

var EmptyHash = Hash{}

type Hash [HashBytes]byte

func NewHash(data []byte) (h Hash) {
	hash := streebog.New()
	n, err := hash.Write(data)
	if err != nil || n != len(data) {
		log.Panicf("calculate hash: %s; n = %d, len(data) = %d", err, n, len(data))
	}

	copy(h[:], hash.Sum(nil))
	return
}

func HashFromSlice(hashData []byte) (hash Hash) {
	if len(hashData) != HashBytes {
		log.Panicf("invalid hash length: %d; %s", len(hashData), string(hashData))
	}

	copy(hash[:], hashData)
	return
}

func HashFromString(str string) (hash Hash) {
	return HashFromSlice(base58.Decode(str))
}

func (h Hash) String() string {
	return base58.Encode(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) isEmpty() bool {
	return h == EmptyHash
}
