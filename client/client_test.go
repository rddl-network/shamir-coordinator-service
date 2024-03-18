package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/client"
	"github.com/stretchr/testify/assert"
)

func TestGetMnemonics(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"mnemonics":["word1","word2"]}`)) // mock response
	}))
	defer mockServer.Close()

	c := client.NewShamirCoordinatorClient(mockServer.URL, mockServer.Client())
	res, err := c.GetMnemonics(context.Background())

	assert.NoError(t, err)
	assert.Len(t, res.Mnemonics, 2)
	assert.Equal(t, "word1", res.Mnemonics[0])
	assert.Equal(t, "word2", res.Mnemonics[1])
}
