package service

import (
	"fmt"
	"strings"
)

func (s *ShamirCoordinatorService) CollectMnemonics() ([]string, error) {
	mnemonics := []string{}
	shareHolderURIs := strings.Split(s.cfg.ShareHolderList, ",")
	for _, shareHolderURI := range shareHolderURIs {
		mnemonic, err := (s.ssc).GetMnemonic(shareHolderURI)
		if err != nil {
			fmt.Printf("Error collecting a share from %s: %s\n", shareHolderURI, err.Error())
		}
		mnemonics = append(mnemonics, mnemonic)
	}
	return mnemonics, nil
}

func (s *ShamirCoordinatorService) deployMnemonics(mnemonics []string) (err error) {
	fmt.Println("ShareHolderUri: " + s.cfg.ShareHolderList)
	shareHolderURIs := strings.Split(s.cfg.ShareHolderList, ",")
	if len(shareHolderURIs) != len(mnemonics) {
		fmt.Println("Error: the amount of shareholders does not match the amount of mnemonics to be deployed: %i shareholders : %i mnemonics",
			len(shareHolderURIs), len(mnemonics))
	}
	for index, shareHolderURI := range shareHolderURIs {
		fmt.Println("ShareHolderUri: " + shareHolderURI)
		err = (s.ssc).PostMnemonic(shareHolderURI, mnemonics[index])
		if err != nil {
			fmt.Printf("Error deploying the sahres at index %d, shareholder %s: %s\n", index, shareHolderURI, err.Error())
			fmt.Println("Attention: redeploy share as there is most likely a inconsistent state")
			return
		}
	}
	return
}
