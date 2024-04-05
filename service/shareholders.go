package service

import (
	"context"
	"fmt"
	"log"
)

func (s *ShamirCoordinatorService) CollectMnemonics() (mnemonics []string, err error) {
	for host, client := range s.sscs {
		s.logger.Info("collecting mnemonic from " + host)
		resp, err := client.GetMnemonic(context.Background())
		if err != nil {
			log.Printf("Error collecting a share from %s: %s\n", host, err.Error())
		} else {
			s.logger.Info("successfully collected")
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
		s.logger.Error(msg)
		return fmt.Errorf(msg)

	}

	i := 0
	for host, client := range s.sscs {
		s.logger.Info("Deploying mnemonic to " + host)
		err = client.PostMnemonic(context.Background(), mnemonics[i])
		if err != nil {
			s.logger.Error("Error deploying the shares at host %s: %s\n", host, err.Error())
			s.logger.Error("Attention: redeploy share as there is most likely a inconsistent state")
			return
		} else {
			s.logger.Info("successfully deployed")
		}
		i++
	}
	return
}
