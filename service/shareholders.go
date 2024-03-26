package service

import (
	"context"
	"fmt"
	"log"
)

func (s *ShamirCoordinatorService) CollectMnemonics() (mnemonics []string, err error) {
	for host, client := range s.sscs {
		resp, err := client.GetMnemonic(context.Background())
		if err != nil {
			log.Printf("Error collecting a share from %s: %s\n", host, err.Error())
		}
		mnemonics = append(mnemonics, resp.Mnemonic)
	}
	return
}

func (s *ShamirCoordinatorService) deployMnemonics(mnemonics []string) (err error) {
	if len(mnemonics) != len(s.sscs) {
		return fmt.Errorf("error: the amount of shareholders does not match the amount of mnemonics to be deployed: %d shareholders : %d mnemonics",
			len(s.sscs), len(mnemonics),
		)
	}

	i := 0
	for host, client := range s.sscs {
		err = client.PostMnemonic(context.Background(), mnemonics[i])
		if err != nil {
			log.Printf("Error deploying the shares at host %s: %s\n", host, err.Error())
			log.Println("Attention: redeploy share as there is most likely a inconsistent state")
			return
		}
		i++
	}
	return
}
