package service_test

import (
	"bytes"
	"context"
	"encoding/json"
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
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(&mycfg, ssc, slip39mock)

	routes := s.GetRoutes()
	assert.Equal(t, 3, len(routes))
}

func TestNotTestMode(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)

	routes := s.GetRoutes()
	assert.Equal(t, 2, len(routes))
}

func TestSendPass(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)

	elements.Client = &elementsmocks.MockClient{}
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)

	request := service.SendTokensRequest{Amount: 123.456, Recipient: "1111111111111111111111111111"}
	jsonString, err := json.Marshal(request)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(jsonString))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"tx-id\":\"0000000000000000000000000000000000000000000000000000000000000000\"}", w.Body.String())
	assert.Equal(t, 200, w.Code)
}

func TestSendFail(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)

	elements.Client = &elementsmocks.MockClient{}
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBufferString("testobject"))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "error")
	assert.Equal(t, 400, w.Code)
}

func TestDeployCheckHex(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)

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
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)

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
	slip39mock := &service.Slip39Interface{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)

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
