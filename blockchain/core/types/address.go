package types

import (
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

const (
	AddressLength        = 20
	InitBalance   uint64 = 1024
)

type Address [AddressLength]byte

var ZeroAddress = Address{}

func AddressFromBytes(data []byte) (address Address) {
	copy(address[:], data)
	return
}

func AddressFromHex(hexData string) (address Address) {
	if len(hexData) != 40 && len(hexData) != 42 {
		log.Panicf("invalid address length: %s", hexData)
	}

	if len(hexData) == 42 {
		hexData = hexData[2:]
	}

	data, err := hex.DecodeString(hexData)
	if err != nil {
		log.Panicf("invalid hex data: %s", hexData)
	}

	copy(address[:], data)
	return
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) Hex() string {
	return hex.EncodeToString(a.Bytes())
}
