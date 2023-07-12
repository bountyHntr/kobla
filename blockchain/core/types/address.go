package types

const AddressLength = 20

type Address [AddressLength]byte

func AddressFromBytes(data []byte) (address Address) {
	copy(address[:], data)
	return
}

func (a Address) Bytes() []byte {
	return a[:]
}
