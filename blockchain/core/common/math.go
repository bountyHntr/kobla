package common

import (
	"encoding/binary"
	"math/big"
)

var (
	BigZero = big.NewInt(0)
	BigOne  = big.NewInt(1)
)

func Int64ToBytes[T int64 | uint64](value T) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(value))
	return data
}

func Int64FromBytes[T int64 | uint64](data []byte) T {
	value := binary.BigEndian.Uint64(data)
	return T(value)
}
