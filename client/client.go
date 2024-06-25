package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rddl-network/shamir-coordinator-service/types"
)

type IShamirCoordinatorClient interface {
	GetMnemonics(ctx context.Context) (res types.MnemonicsResponse, err error)
	PostMnemonics(ctx context.Context, secret string) (err error)
	SendTokens(ctx context.Context, recipient string, amount string) (res types.SendTokensResponse, err error)
}

type ShamirCoordinatorClient struct {
	baseURL string
	client  *http.Client
}

func NewShamirCoordinatorClient(baseURL string, client *http.Client) *ShamirCoordinatorClient {
	if client == nil {
		client = &http.Client{}
	}
	return &ShamirCoordinatorClient{
		baseURL: baseURL,
		client:  client,
	}
}

func (scc *ShamirCoordinatorClient) GetMnemonics(ctx context.Context) (res types.MnemonicsResponse, err error) {
	err = scc.doRequest(ctx, http.MethodGet, scc.baseURL+"/mnemonics", nil, &res)
	return
}

func (scc *ShamirCoordinatorClient) PostMnemonics(ctx context.Context, secret string) (err error) {
	err = scc.doRequest(ctx, http.MethodPost, scc.baseURL+"/mnemonics/"+url.PathEscape(secret), nil, nil)
	return
}

func (scc *ShamirCoordinatorClient) SendTokens(ctx context.Context, recipient string, amount string) (res types.SendTokensResponse, err error) {
	requestBody := types.SendTokensRequest{
		Recipient: recipient,
		Amount:    amount,
	}
	err = scc.doRequest(ctx, http.MethodPost, scc.baseURL+"/send", &requestBody, &res)
	return
}

func (scc *ShamirCoordinatorClient) ReIssueAsset(ctx context.Context, asset string, amount string) (res types.ReIssueResponse, err error) {
	requestBody := types.ReIssueRequest{
		Asset:  asset,
		Amount: amount,
	}
	err = scc.doRequest(ctx, http.MethodPost, scc.baseURL+"/reissue", &requestBody, &res)
	return
}

func (scc *ShamirCoordinatorClient) doRequest(ctx context.Context, method, url string, body interface{}, response interface{}) (err error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := scc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &httpError{StatusCode: resp.StatusCode, Msg: strings.Join(resp.Header["Error"], "\n")}
	}

	if response != nil {
		return json.NewDecoder(resp.Body).Decode(response)
	}

	return
}

type httpError struct {
	StatusCode int
	Msg        string
}

func (e *httpError) Error() string {
	return http.StatusText(e.StatusCode) + ": " + e.Msg
}
