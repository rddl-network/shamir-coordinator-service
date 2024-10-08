package service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	elements "github.com/rddl-network/elements-rpc"
	elementsmocks "github.com/rddl-network/elements-rpc/utils/mocks"
	log "github.com/rddl-network/go-utils/logger"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/rddl-network/shamir-coordinator-service/types"
	"github.com/rddl-network/shamir-shareholder-service/client"
	"github.com/stretchr/testify/assert"
)

func TestTestMode(t *testing.T) {
	cfg := config.GetConfig()
	mycfg := *cfg
	mycfg.TestMode = true

	sscs := make(map[string]client.IShamirShareholderClient)
	ctrl := gomock.NewController(t)
	ssc := testutil.NewMockIShamirShareholderClient(ctrl)
	sscs["client"] = ssc

	slip39mock := &testutil.Slip39Mock{}
	logger := log.GetLogger(log.DEBUG)
	db := testutil.SetupTestDBConnector(t)
	s := service.NewShamirCoordinatorService(&mycfg, sscs, slip39mock, logger, db)

	routes := s.GetRoutes()
	assert.Equal(t, 5, len(routes))
}

func TestNotTestMode(t *testing.T) {
	s := testutil.SetupTestService(t)

	routes := s.GetRoutes()
	assert.Equal(t, 4, len(routes))
}

func TestSendPass(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	request := types.SendTokensRequest{Amount: "123.456", Recipient: "1111111111111111111111111111"}
	jsonString, err := json.Marshal(request)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(jsonString))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"tx-id\":\"0000000000000000000000000000000000000000000000000000000000000000\"}", w.Body.String())
	assert.Equal(t, 200, w.Code)
}

func TestSendFailAmount0(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	request := types.SendTokensRequest{Amount: "0.0", Recipient: "1111111111111111111111111111"}
	jsonString, err := json.Marshal(request)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(jsonString))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestSendFail(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBufferString("testobject"))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "Error")
	assert.Equal(t, 400, w.Code)
}

func TestDeployCheckHex(t *testing.T) {
	s := testutil.SetupTestService(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonics/öaksjdf", nil)
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"Error\":\"the secret has to be send in valid hex string format\"}", w.Body.String())
	assert.Equal(t, 500, w.Code)
}

func TestDeployCheckLength(t *testing.T) {
	s := testutil.SetupTestService(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonics/abcdef", nil)
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"Error\":\"the secret has to be of length 32 or 64 (16 or 32 byte)\"}", w.Body.String())
	assert.Equal(t, 500, w.Code)

	w = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, "/mnemonics/abcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef", nil)
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"Error\":\"the secret has to be of length 32 or 64 (16 or 32 byte)\"}", w.Body.String())
	assert.Equal(t, 500, w.Code)
}

func TestDeployPass(t *testing.T) {
	s := testutil.SetupTestServiceWithSlip39Interface(t)

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

func TestReissuePass(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	request := types.ReIssueRequest{Amount: "123.456", Asset: "06c20c8de513527f1ae6c901f74a05126525ac2d7e89306f4a7fd5ec4e674403"}
	jsonString, err := json.Marshal(request)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/reissue", bytes.NewBuffer(jsonString))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"tx-id\":\"0000000000000000000000000000000000000000000000000000000000000000\"}", w.Body.String())
	assert.Equal(t, 200, w.Code)
}

func TestReissueFail(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/reissue", bytes.NewBufferString("testobject"))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "Error")
	assert.Equal(t, 400, w.Code)
}

func TestIssueMachineNFTPass(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	request := types.IssueMachineNFTRequest{Name: "Machine", MachineAddress: "someAddr", Domain: "testnet-assets.rddl.io"}
	jsonString, err := json.Marshal(request)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/issue-machine-nft", bytes.NewBuffer(jsonString))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Equal(t, "{\"asset\":\"0000000000000000000000000000000000000000000000000000000000000000\",\"contract\":\"{\\\"entity\\\":{\\\"domain\\\":\\\"testnet-assets.rddl.io\\\"},\\\"issuer_pubkey\\\":\\\"020000000000000000000000000000000000000000000000000000000000000000\\\",\\\"machine_addr\\\":\\\"someAddr\\\",\\\"name\\\":\\\"Machine\\\",\\\"precision\\\":0,\\\"version\\\":0}\",\"hex-tx\":\"0000000000000000000000000000000000000000000000000000000000000000\"}", w.Body.String())
	assert.Equal(t, 200, w.Code)
}

func TestIssueMachineNFTFail(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/issue-machine-nft", bytes.NewBufferString("testobject"))
	assert.NoError(t, err)
	s.Router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "Error")
	assert.Equal(t, 400, w.Code)
}

type BitcoinAddressError struct {
	msg string
}

func (e BitcoinAddressError) Error() string {
	return e.msg
}

func TestAddToQueue(t *testing.T) {
	s := testutil.SetupTestService(t)
	err := BitcoinAddressError{"Invalid Bitcoin Address:"}
	assert.False(t, s.AddToQueue(err))

	err = BitcoinAddressError{": -5"}
	assert.False(t, s.AddToQueue(err))
}
