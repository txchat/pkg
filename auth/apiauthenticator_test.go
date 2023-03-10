package auth

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	xcrypt "github.com/txchat/pkg/crypt"
	secp256k1_ethereum "github.com/txchat/pkg/crypt/secp256k1-ethereum"
	secp256k1_haltingstate "github.com/txchat/pkg/crypt/secp256k1-haltingstate"
)

var userAddress = "1AqutxNoVTtcWiVYpBtvficAgea1dYTddR"

func TestDefaultAuthAndVerify(t *testing.T) {
	pubKey, err := hex.DecodeString(publicKey)
	assert.Nil(t, err)

	privKey, err := hex.DecodeString(privateKey)
	assert.Nil(t, err)

	client := NewDefaultAPIAuthenticator()
	sig := client.Request("dtalk", pubKey, privKey)

	server := NewDefaultAPIAuthenticator()
	uid, err := server.Auth(sig)
	assert.Nil(t, err)
	assert.Equal(t, userAddress, uid)
}

func TestHaltAuthAndEthVerify(t *testing.T) {
	pubKey, err := hex.DecodeString(publicKey)
	assert.Nil(t, err)

	privKey, err := hex.DecodeString(privateKey)
	assert.Nil(t, err)

	haltDriver, err := xcrypt.Load(secp256k1_haltingstate.Name)
	if err != nil {
		panic(err)
	}
	ethereumDriver, err := xcrypt.Load(secp256k1_ethereum.Name)
	if err != nil {
		panic(err)
	}
	client := NewDefaultAPIAuthenticatorAsDriver(haltDriver)
	sig := client.Request("dtalk", pubKey, privKey)

	server := NewDefaultAPIAuthenticatorAsDriver(ethereumDriver)
	uid, err := server.Auth(sig)
	assert.Nil(t, err)
	assert.Equal(t, userAddress, uid)
}

func TestEthAuthAndHaltVerify(t *testing.T) {
	pubKey, err := hex.DecodeString(publicKey)
	assert.Nil(t, err)

	privKey, err := hex.DecodeString(privateKey)
	assert.Nil(t, err)

	haltDriver, err := xcrypt.Load(secp256k1_haltingstate.Name)
	if err != nil {
		panic(err)
	}
	ethereumDriver, err := xcrypt.Load(secp256k1_ethereum.Name)
	if err != nil {
		panic(err)
	}
	client := NewDefaultAPIAuthenticatorAsDriver(ethereumDriver)
	sig := client.Request("dtalk", pubKey, privKey)

	server := NewDefaultAPIAuthenticatorAsDriver(haltDriver)
	uid, err := server.Auth(sig)
	assert.Nil(t, err)
	assert.Equal(t, userAddress, uid)
}

func TestDefaultVerifyInvalid(t *testing.T) {
	sig := "#1655716335340*dtalk#022db0e08669b30c5dab8c564b428db4944912144088943ec9b690a9046bc8f78b"
	server := NewDefaultAPIAuthenticator()
	uid, err := server.Auth(sig)
	assert.ErrorAs(t, err, &SignatureInvalidError{})
	assert.Empty(t, uid)
}

func TestDefaultAuthExpire(t *testing.T) {
	sig := "zrrLQ9FLnpON9s3erkMJ+sug5oviPYcOR04/w4ucC5dVbkjqvuZIFUuGgaKi+5XmDzj2FvxrDIom9dt2Ons6fAE=#1655716335340*dtalk#022db0e08669b30c5dab8c564b428db4944912144088943ec9b690a9046bc8f78b"
	server := NewDefaultAPIAuthenticator()
	uid, err := server.Auth(sig)
	assert.ErrorIs(t, err, ErrSignatureExpired)
	assert.Empty(t, uid)
}

func TestError(t *testing.T) {
	sig := "#1655716335340*dtalk#022db0e08669b30c5dab8c564b428db4944912144088943ec9b690a9046bc8f78b"
	server := NewDefaultAPIAuthenticator()
	uid, err := server.Auth(sig)
	assert.Empty(t, uid)

	assert.True(t, errors.As(err, &SignatureInvalidError{}))
	assert.ErrorAs(t, err, &SignatureInvalidError{})
	// once unWarp not nil
	assert.NotNil(t, errors.Unwrap(err))
	// twice unWarp is nil
	assert.Nil(t, errors.Unwrap(errors.Unwrap(err)))
}
