package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rddl-network/shamir-coordinator-service/config"
)

type IShamirShareholderClient interface {
	GetMnemonic(shareholderURI string) (string, error)
	PostMnemonic(shareholderURI string, mnemonic string) error
}

type ShamirShareholderClient struct {
	cfg *config.Config
}

func NewShamirShareholderClient(cfg *config.Config) *ShamirShareholderClient {
	ssc := &ShamirShareholderClient{}
	ssc.cfg = cfg
	return ssc
}

type ShareHolderResponse struct {
	Mnemonic string `json:"mnemonic"`
}

func (s *ShamirShareholderClient) GetMnemonic(shareholderURI string) (mnemonic string, err error) {
	client, err := s.get2wayTLSClient()
	if err != nil {
		fmt.Printf("Error creating the 2WayTLS client: %s\n", err.Error())
		return
	}

	// Make request
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, shareholderURI+"/mnemonic", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	jsonBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// The struct that matches the JSON structure
	var response ShareHolderResponse

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(jsonBody, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", jsonBody)
	mnemonic = response.Mnemonic
	return
}

func (s *ShamirShareholderClient) PostMnemonic(shareHolderURI string, mnemonic string) (err error) {
	client, err := s.get2wayTLSClient()
	if err != nil {
		fmt.Printf("Error creating the 2WayTLS client: %s\n", err.Error())
		return
	}

	jsonString := fmt.Sprintf(`{"mnemonic":"%s"}`, mnemonic)
	jsonData := []byte(jsonString)

	// Create new request with POST method and JSON data
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, shareHolderURI+"/mnemonic", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error performing request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response body
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	return
}

func (s *ShamirShareholderClient) get2wayTLSClient() (client *http.Client, err error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(s.cfg.CertsPath+"client.crt", s.cfg.CertsPath+"client.key")
	if err != nil {
		fmt.Printf("Error loading client certificate: %v\n", err)
		return
	}

	// Load CA cert
	caCert, err := os.ReadFile(s.cfg.CertsPath + "ca.crt")
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
