package service

import (
	"encoding/hex"
	"log"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	hexutil "github.com/rddl-network/go-utils/hex"
)

func getRandomPrivateKey(n int) (string, error) {
	return hexutil.RandomHex(n)
}

func GenerateNewKeyPair() (*secp256k1.PrivateKey, *secp256k1.PublicKey) {
	pkSource, _ := getRandomPrivateKey(32)
	privateKeyBytes, err := hex.DecodeString(pkSource)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// Initialize a secp256k1 private key object.
	privateKey, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)
	return privateKey, publicKey
}
