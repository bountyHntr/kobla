package types_test

import (
	"crypto/rand"
	"testing"

	"kobla/blockchain/core/types"

	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	data := make([]byte, 100)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	hash := types.NewHash(data)
	assert.NotEmpty(t, hash)
	assert.Equal(t, hash[:], hash.Bytes())
	assert.Equal(t, hash, types.HashFromSlice(hash.Bytes()))
	assert.Zero(t, types.EmptyHash)

	hashStr := hash.String()
	assert.NotEmpty(t, hashStr)
	assert.Equal(t, hash.Bytes(), base58.Decode(hashStr))
	assert.Equal(t, hash, types.HashFromString(hashStr))
}
