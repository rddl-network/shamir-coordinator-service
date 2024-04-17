package service

import (
	"context"
	"fmt"
)

func (s *ShamirCoordinatorService) CollectMnemonics() (mnemonics []string, err error) {
	for host, client := range s.sscs {
		resp, err := client.GetMnemonic(context.Background())
		if err != nil {
			msg := fmt.Sprintf("Error collecting a share from %s: %s\n", host, err.Error())
			s.logger.Error("msg", msg)
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
		s.logger.Error("msg", msg)
		return fmt.Errorf(msg)
	}

	i := 0
	for host, client := range s.sscs {
		err = client.PostMnemonic(context.Background(), mnemonics[i])
		if err != nil {
			msg := fmt.Sprintf("Error deploying the shares at host %s: %s\n", host, err.Error())
			s.logger.Error("msg", msg)
			s.logger.Error("msg", "Attention: redeploy share as there is most likely a inconsistent state")
			return
		}
		i++
	}
	return
}
