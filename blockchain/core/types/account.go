package types

import (
	"crypto/rand"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	signf "github.com/ddulesov/gogost/gost3410"
)

const (
	AddressLength        = 128
	InitBalance   uint64 = 10240
)

var (
	curve = signf.CurveIdtc26gost341012512paramSetA()
	mode  = signf.Mode2012
)

type Account struct {
	prv *signf.PrivateKey
	pub *signf.PublicKey
}

func NewAccount() (Account, error) {
	prv, err := signf.GenPrivateKey(curve, mode, rand.Reader)
	if err != nil {
		return Account{}, fmt.Errorf("generate private key: %w", err)
	}

	pub, err := prv.PublicKey()
	if err != nil {
		return Account{}, fmt.Errorf("calculate public key: %w", err)
	}

	return Account{prv, pub}, nil
}

func AccountFromPrivKey(privateKey string) (Account, error) {
	prv, err := signf.NewPrivateKey(curve, mode, base58.Decode(privateKey))
	if err != nil {
		return Account{}, fmt.Errorf("private key: %w", err)
	}

	pub, err := prv.PublicKey()
	if err != nil {
		return Account{}, fmt.Errorf("calculate public key: %w", err)
	}

	return Account{prv, pub}, nil
}

func (a Account) Address() (address Address) {
	copy(address[:], a.pub.Raw())
	return
}

func (a Account) Sign(hash Hash) ([]byte, error) {
	return a.prv.SignDigest(hash[:], rand.Reader)
}

func (a Account) Verify(hash Hash, signature []byte) (bool, error) {
	return a.pub.VerifyDigest(hash[:], signature)
}

func (a Account) PrivateKey() string {
	return base58.Encode(a.prv.Raw())
}

type Address [AddressLength]byte

var ZeroAddress = Address{}

func AddressFromBytes(data []byte) (address Address) {
	copy(address[:], data)
	return
}

func AddressFromString(str string) (address Address) {
	copy(address[:], base58.Decode(str))
	return
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) String() string {
	return base58.Encode(a.Bytes())
}

func (a Address) Verify(hash Hash, signature []byte) (bool, error) {
	pub, err := signf.NewPublicKey(curve, mode, a[:])
	if err != nil {
		return false, fmt.Errorf("public key: %w", err)
	}

	return pub.VerifyDigest(hash[:], signature)
}
