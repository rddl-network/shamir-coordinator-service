package service

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/rddl-network/shamir-coordinator-service/config"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getRandomPrivateKey(n int) (string, error) {
	return randomHex(n)
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

func IsValidHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

func ContainsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func Get2wayTLSClient(cfg *config.Config) (client *http.Client, err error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(cfg.CertsPath+"client.crt", cfg.CertsPath+"client.key")
	if err != nil {
		fmt.Printf("Error loading client certificate: %v\n", err)
		return
	}

	// Load CA cert
	caCert, err := os.ReadFile(cfg.CertsPath + "ca.crt")
	if err != nil {
		fmt.Printf("Error loading CA certificate: %v\n", err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	return
}
