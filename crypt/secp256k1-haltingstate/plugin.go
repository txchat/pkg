package secp256K1

import (
	"github.com/txchat/pkg/crypt"
)

const Name = "secp256k1-haltingstate"

func init() {
	crypt.Register(Name, New())
}

func New() crypt.Encrypt {
	return &haltingstate{}
}
