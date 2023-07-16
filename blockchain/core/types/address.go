package types

import "encoding/hex"

const AddressLength = 20

type Address [AddressLength]byte

var ZeroAddress = Address{}

func AddressFromBytes(data []byte) (address Address) {
	copy(address[:], data)
	return
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) Hex() string {
	return hex.EncodeToString(a.Bytes())
}
