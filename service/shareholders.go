package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func (s *ShamirCoordinatorService) CollectMnemonics() ([]string, error) {
	mnemonics := []string{}
	shareHolderURIs := strings.Split(s.cfg.ShareHolderList, ",")
	for _, shareHolderUri := range shareHolderURIs {
		mnemonic, err := s.collectShare(shareHolderUri)
		if err != nil {
			fmt.Printf("Error collecting a share from %s: %s\n", shareHolderUri, err.Error())
		}
		mnemonics = append(mnemonics, mnemonic)
	}
	return mnemonics, nil
}

func (s *ShamirCoordinatorService) deployMnemonics(mnemonics []string) (err error) {

	fmt.Println("ShareHolderUri: " + s.cfg.ShareHolderList)
	shareHolderURIs := strings.Split(s.cfg.ShareHolderList, ",")
	if len(shareHolderURIs) != len(mnemonics) {
		fmt.Println("Error: the amount of shareholders does not match the amount of mnemonics to be deployed: %i shareholders : %i mnemonics",
			len(shareHolderURIs), len(mnemonics))
	}
	for index, shareHolderUri := range shareHolderURIs {
		fmt.Println("ShareHolderUri: " + shareHolderUri)
		err = s.setShareHolderMnemonic(shareHolderUri, mnemonics[index])
		if err != nil {
			fmt.Printf("Error deploying the sahres at index %d, shareholder %s: %s\n", index, shareHolderUri, err.Error())
			fmt.Println("Attention: redeploy share as there is most likely a inconsistent state")
			return
		}
	}
	return
}

type ShareHolderResponse struct {
	Mnemonic string `json:"mnemonic"`
}

func (s *ShamirCoordinatorService) collectShare(shareholderURI string) (mnemonic string, err error) {
	client, err := s.get2wayTLSClient()
	if err != nil {
		fmt.Printf("Error creating the 2WayTLS client: %s\n", err.Error())
		return
	}

	// Make request
	resp, err := client.Get(shareholderURI + "/mnemonic")
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
	err = json.Unmarshal([]byte(jsonBody), &response)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	fmt.Printf("Response: %s\n", jsonBody)
	mnemonic = response.Mnemonic
	return
}

func (s *ShamirCoordinatorService) get2wayTLSClient() (client *http.Client, err error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(s.cfg.CertsPath+"server.crt", s.cfg.CertsPath+"server.key")
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
	//tlsConfig.BuildNameToCertificate()
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	return
}

func (s *ShamirCoordinatorService) setShareHolderMnemonic(shareHolderUri string, mnemonic string) (err error) {
	client, err := s.get2wayTLSClient()
	if err != nil {
		fmt.Printf("Error creating the 2WayTLS client: %s\n", err.Error())
		return
	}

	jsonString := fmt.Sprintf(`{"mnemonic":"%s"}`, mnemonic)
	jsonData := []byte(jsonString)

	// Create new request with POST method and JSON data
	req, err := http.NewRequest("POST", shareHolderUri+"/mnemonic", bytes.NewBuffer(jsonData))
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
