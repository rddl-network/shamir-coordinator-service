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
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"mnemonics":["word1","word2"]}`))
	}))
	defer mockServer.Close()

	c := client.NewShamirCoordinatorClient(mockServer.URL, mockServer.Client())
	res, err := c.GetMnemonics(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res.Mnemonics, 2)
	assert.Equal(t, "word1", res.Mnemonics[0])
	assert.Equal(t, "word2", res.Mnemonics[1])
}

func TestPostMnemonics(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/mnemonics/someSecret", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	c := client.NewShamirCoordinatorClient(mockServer.URL, mockServer.Client())
	err := c.PostMnemonics(context.Background(), "someSecret")

	assert.NoError(t, err)
}

func TestSendTokens(t *testing.T) {
	expectedRequestBody := `{"recipient":"testRecipient","amount":"123"}`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, expectedRequestBody, string(bodyBytes))
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/send", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"tx-id":"12345"}`))
	}))
	defer mockServer.Close()

	c := client.NewShamirCoordinatorClient(mockServer.URL, mockServer.Client())
	res, err := c.SendTokens(context.Background(), "testRecipient", "123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "12345", res.TxID)
}
