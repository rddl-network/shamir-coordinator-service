package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/rddl-network/shamir-coordinator-service/service"
)

type IShamirCoordinatorClient interface {
	GetMnemonics() (res service.MnemonicsResponse, err error)
	PostMnemonics(secret string) (err error)
	SendTokens(recipient string, amount string) (res service.SendTokensResponse, err error)
}

type ShamirCoordinatorClient struct {
	host   string
	client *http.Client
}

func NewShamirCoordinatorClient(host string) *ShamirCoordinatorClient {
	return &ShamirCoordinatorClient{
		host:   host,
		client: &http.Client{},
	}
}

func (scc *ShamirCoordinatorClient) GetMnemonics() (res service.MnemonicsResponse, err error) {
	url := &url.URL{
		Scheme: "https",
		Host:   scc.host,
		Path:   "/mnemonics",
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url.String(), nil)
	if err != nil {
		return
	}

	resp, err := scc.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyBytes, &res)
	return
}

func (scc *ShamirCoordinatorClient) PostMnemonics(secret string) (err error) {
	url := &url.URL{
		Scheme: "https",
		Host:   scc.host,
		Path:   "/mnemonics/" + secret,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url.String(), nil)
	if err != nil {
		return
	}

	resp, err := scc.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}

func (scc *ShamirCoordinatorClient) SendTokens(recipeint string, amount string) (res *service.SendTokensResponse, err error) {
	url := &url.URL{
		Scheme: "https",
		Host:   scc.host,
		Path:   "/send",
	}

	body := &service.SendTokensRequest{
		Recipient: recipeint,
		Amount:    amount,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return
	}

	resp, err := scc.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyBytes, &res)
	return
}
