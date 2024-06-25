package service

import (
	"context"
	"errors"
	"fmt"
)

func (s *ShamirCoordinatorService) CollectMnemonics() (mnemonics []string, err error) {
	for host, client := range s.sscs {
		s.logger.Info("msg", "collecting mnemonic from "+host)
		resp, err := client.GetMnemonic(context.Background())
		if err != nil {
			s.logger.Error("error", "Error collecting a share from "+host+" "+err.Error())
		} else {
			s.logger.Info("msg", "successfully collected")
		}
		mnemonics = append(mnemonics, resp.Mnemonic)
	}
	return
}

func (s *ShamirCoordinatorService) deployMnemonics(mnemonics []string) (err error) {
	if len(mnemonics) != len(s.sscs) {
		msg := fmt.Sprintf("error: the amount of shareholders does not match the amount of mnemonics to be deployed: %d shareholders : %d mnemonics",
			len(s.sscs), len(mnemonics),
		)
		s.logger.Error("error", msg)
		err = errors.New(msg)
		return
	}

	i := 0
	for host, client := range s.sscs {
		s.logger.Info("msg", "Deploying mnemonic to "+host)
		err = client.PostMnemonic(context.Background(), mnemonics[i])
		if err != nil {
			s.logger.Error("error", "Error deploying the shares at host "+host+" "+err.Error())
			s.logger.Error("error", "Attention: redeploy share as there is most likely a inconsistent state")
			return
		}
		s.logger.Info("msg", "successfully deployed")
		i++
	}
	return
}
