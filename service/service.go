package service

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/rddl-network/go-utils/logger"
	"github.com/rddl-network/go-utils/tls"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service/backend"
	"github.com/rddl-network/shamir-shareholder-service/client"
)

type ShamirCoordinatorService struct {
	cfg             *config.Config
	Router          *gin.Engine
	sscs            map[string]client.IShamirShareholderClient
	slip39Interface ISlip39
	logger          log.AppLogger
	db              *backend.DBConnector
}

func NewShamirCoordinatorService(cfg *config.Config, sscs map[string]client.IShamirShareholderClient, slip39Interface ISlip39, logger log.AppLogger, db *backend.DBConnector) *ShamirCoordinatorService {
	service := &ShamirCoordinatorService{
		cfg:             cfg,
		sscs:            sscs,
		slip39Interface: slip39Interface,
		logger:          logger,
		db:              db,
	}

	gin.SetMode(gin.ReleaseMode)
	service.Router = gin.New()
	service.Router.POST("/send", service.SendTokens)
	service.Router.POST("/reissue", service.ReIssue)
	service.Router.POST("/issue-machine-nft", service.IssueMachineNFT)
	service.Router.POST("/mnemonics/:secret", service.DeployShares)
	if cfg.TestMode {
		service.Router.GET("/mnemonics", service.CollectShares)
	}
	return service
}

func (s *ShamirCoordinatorService) Run() (err error) {
	cfg := config.GetConfig()
	caCertFile, err := os.ReadFile(cfg.CertsPath + "ca.crt")
	if err != nil {
		return err
	}

	tlsConfig := tls.Get2WayTLSServer(caCertFile)
	server := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", cfg.ServiceBind, cfg.ServicePort),
		TLSConfig: tlsConfig,
		Handler:   s.Router,
	}

	// workaround to listen on tcp4 and not tcp6
	// https://stackoverflow.com/a/38592286
	ln, err := net.Listen("tcp4", server.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	go s.rerunFailedRequests(cfg.WaitPeriod)
	s.logger.Info("msg", "started server", "host", cfg.ServiceBind, "port", cfg.ServicePort)
	return server.ServeTLS(ln, cfg.CertsPath+"server.crt", cfg.CertsPath+"server.key")
}

func (s *ShamirCoordinatorService) rerunFailedRequests(waitPeriod int) {
	ticker := time.NewTicker(time.Duration(waitPeriod) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		passphrase, err := s.GetPassphrase()
		if err != nil {
			s.logger.Error("error", errWalletMsg+err.Error())
			continue
		}

		// prepare the wallet, loading and unlocking
		err = s.PrepareWallet(passphrase)
		if err != nil {
			s.logger.Error("error", errWalletMsg+err.Error())
			continue
		}

		sendTokensRequests, err := s.db.GetAllSendTokensRequests()
		if err != nil {
			s.logger.Error("msg", "error while reading sendTokensRequests: "+err.Error())
		}
		for _, req := range sendTokensRequests {
			s.handleSendTokensRequest(req)
		}

		reIssueRequests, err := s.db.GetAllReissueRequests()
		if err != nil {
			s.logger.Error("msg", "error while reading reIssueRequests: "+err.Error())
		}
		for _, req := range reIssueRequests {
			s.handleReIssueRequest(req)
		}

		issueNFTAssetRequests, err := s.db.GetAllIssueMachineNFTRequests()
		if err != nil {
			s.logger.Error("msg", "error while reading issueNFTAssetRequests: "+err.Error())
		}
		for _, req := range issueNFTAssetRequests {
			s.handleIssueMachineNFTRequest(req)
		}

		if err = s.WalletLock(); err != nil {
			s.logger.Error("error", errWalletLockMsg+err.Error())
		}
	}
}

func (s *ShamirCoordinatorService) handleSendTokensRequest(req backend.SendTokensRequest) {
	txID, err := s.SendAsset(req.Recipient, req.Amount, req.Asset)
	if err != nil {
		s.logger.Error("error", "error sending the transaction: "+err.Error())
		return
	}
	s.logger.Info("msg", "successfully sended tx with id: "+txID+" to "+req.Recipient)
	if err = s.db.DeleteRequest(backend.SendTokensRequestPrefix, req.ID); err != nil {
		s.logger.Error("error", "failed to delete SendTokensRequest", "id", req.ID)
	}
}

func (s *ShamirCoordinatorService) handleReIssueRequest(req backend.ReIssueRequest) {
	txID, err := s.ReissueAsset(req.Asset, req.Amount)
	if err != nil {
		s.logger.Error("error", "error reissuing asset: "+err.Error())
		return
	}
	s.logger.Info("msg", "successfully reissued asset", "tx-id", txID, "asset", req.Asset, "amount", req.Amount)
	if err = s.db.DeleteRequest(backend.ReissueRequestPrefix, req.ID); err != nil {
		s.logger.Error("error", "failed to delete ReIssueRequest", "id", req.ID)
	}
}

func (s *ShamirCoordinatorService) handleIssueMachineNFTRequest(req backend.IssueMachineNFTRequest) {
	asset, contract, hexTx, err := s.IssueNFTAsset(req.Name, req.MachineAddress, req.Domain)
	if err != nil {
		s.logger.Error("error", "error issuing machine nft: "+err.Error(), "name", req.Name, "machineAddress", req.MachineAddress, "domain", req.Domain)
		return
	}
	s.logger.Info("msg", "successfully issued machine nft", "asset_id", asset, "contract", contract, "hex_tx", hexTx)
	if err = s.db.DeleteRequest(backend.IssueMachineNFTPrefix, req.ID); err != nil {
		s.logger.Error("error", "failed to delete IssueMachineNFTRequest", "id", req.ID)
	}
}
