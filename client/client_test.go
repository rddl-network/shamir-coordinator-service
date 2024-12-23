package client_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/client"
	"github.com/stretchr/testify/assert"
)

func TestGetMnemonics(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mnemonics", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"mnemonics":["word1","word2"]}`))
		assert.NoError(t, err)
	}))
	defer mockServer.Close()

	c := client.NewSCClient(mockServer.URL, mockServer.Client())
	res, err := c.GetMnemonics(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res.Mnemonics, 2)
	assert.Equal(t, "word1", res.Mnemonics[0])
	assert.Equal(t, "word2", res.Mnemonics[1])
}

func TestPostMnemonics(t *testing.T) {
	t.Parallel()
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mnemonics/someSecret", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	c := client.NewSCClient(mockServer.URL, mockServer.Client())
	err := c.PostMnemonics(context.Background(), "someSecret")

	assert.NoError(t, err)
}

func TestSendTokens(t *testing.T) {
	t.Parallel()
	expectedRequestBody := `{"recipient":"testRecipient","amount":"123","asset":""}`
	expectedResponseBody := `{"tx-id":"12345"}`

	mockServer := setupMockServer(t, http.MethodPost, "/send", expectedRequestBody, expectedResponseBody)
	defer mockServer.Close()

	c := client.NewSCClient(mockServer.URL, mockServer.Client())
	res, err := c.SendTokens(context.Background(), "testRecipient", "123", "")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "12345", res.TxID)
}

func TestReissueAsset(t *testing.T) {
	t.Parallel()
	expectedRequestBody := `{"asset":"06c20c8de513527f1ae6c901f74a05126525ac2d7e89306f4a7fd5ec4e674403","amount":"123"}`
	expectedResponseBody := `{"tx-id":"12345"}`

	mockServer := setupMockServer(t, http.MethodPost, "/reissue", expectedRequestBody, expectedResponseBody)
	defer mockServer.Close()

	c := client.NewSCClient(mockServer.URL, mockServer.Client())
	res, err := c.ReIssueAsset(context.Background(), "06c20c8de513527f1ae6c901f74a05126525ac2d7e89306f4a7fd5ec4e674403", "123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "12345", res.TxID)
}

func TestIssueMachineNFT(t *testing.T) {
	t.Parallel()
	expectedRequestBody := `{"name":"Machine","machine-address":"someAddr","domain":"testnet-assets.rddl.io"}`
	expectedResponseBody := `{"asset":"12345","contract":"someContract","hex-tx":"06c20c8de513527f1ae6c901f74a05126525ac2d7e89306f4a7fd5ec4e674403"}`

	mockServer := setupMockServer(t, http.MethodPost, "/issue-machine-nft", expectedRequestBody, expectedResponseBody)
	defer mockServer.Close()

	c := client.NewSCClient(mockServer.URL, mockServer.Client())
	res, err := c.IssueMachineNFT(context.Background(), "Machine", "someAddr", "testnet-assets.rddl.io")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "12345", res.Asset)
	assert.Equal(t, "someContract", res.Contract)
	assert.Equal(t, "06c20c8de513527f1ae6c901f74a05126525ac2d7e89306f4a7fd5ec4e674403", res.HexTX)
}

func setupMockServer(t *testing.T, method string, route string, expectedRequestBody string, expectedResponseBody string) (mockServer *httptest.Server) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, expectedRequestBody, string(bodyBytes))
		assert.Equal(t, method, r.Method)
		assert.Equal(t, route, r.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(expectedResponseBody))
		assert.NoError(t, err)
	}))
	return
}
