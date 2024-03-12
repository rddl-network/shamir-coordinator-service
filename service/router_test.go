package service_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	elements "github.com/rddl-network/elements-rpc"
	elementsmocks "github.com/rddl-network/elements-rpc/utils/mocks"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTestMode(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	mycfg := *cfg
	mycfg.TestMode = true

	ssc := testutil.NewShamirShareholderClientMock(&mycfg)
	s := service.NewShamirCoordinatorService(&mycfg, ssc)

	routes := s.GetRoutes()
	assert.Equal(t, 3, len(routes))
}

func TestNotTestMode(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	routes := s.GetRoutes()
	assert.Equal(t, 2, len(routes))
}

func TestSendTo(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}

	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)
	address := "tlq1qqvsmfp0w3dmvwtkfteanzk0n7wksu6zx4pywzvak9p6d34yngghw39ynqwcxqrq3muqxffflmprr9exn8ldm79mlkz7dmpy0e"
	amount := "0.0001"
	txID, err := s.SendAsset(address, amount)
	assert.NoError(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", txID)
}

func TestDeployCheckHex(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonics/öaksjdf", nil)
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"error\":\"the secret has to be send in valid hex string format\"}", w.Body.String())
	assert.Equal(t, 500, w.Code)
}

func TestDeployCheckLength(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonics/abcdef", nil)
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"error\":\"the secret has to be of length 32 or 64 (16 or 32 byte)\"}", w.Body.String())
	assert.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonics/abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef", nil)
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"error\":\"the secret has to be of length 32 or 64 (16 or 32 byte)\"}", w.Body.String())
	assert.Equal(t, 500, w.Code)
}

func TestDeployPass(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonics/abcdefabcdefabcdefabcdefabcdef23", nil)
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonics/abcdefabcdefabcdefabcdefabcdef23abcdefabcdefabcdefabcdefabcdef23", nil)
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
