package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	slip39 "github.com/rddl-network/bc-slip39-go"
)

var password = ""

func (s *ShamirCoordinatorService) CreateMnemonics(hexSecret string) (mnemonics []string, err error) {
	// Define the command and arguments
	groupThreshold := uint8(1)
	groups := []slip39.GroupDescriptor{
		{
			Threshold: uint8(s.cfg.ShamirThreshold),
			Count:     uint8(s.cfg.ShamirShares),
		},
	}
	secret, err := hex.DecodeString(hexSecret)
	if err != nil {
		return
	}
	iterationExponent := uint8(0)
	count, wordsInEachShare, sharesBuffer, err := s.slip39Interface.Generate(groupThreshold, groups, secret, password, iterationExponent, slip39.Random())
	if err != nil {
		return
	}

	mnemonics = make([]string, count)
	for index := 0; index < count; index++ {
		start := index * wordsInEachShare
		end := start + wordsInEachShare
		words := sharesBuffer[start:end]
		resultString, err := s.slip39Interface.StringsForWords(words, wordsInEachShare)
		if err != nil {
			return nil, err
		}
		mnemonics[index] = resultString
	}

	if len(mnemonics) != s.cfg.ShamirShares {
		msg := fmt.Sprintf("wrong amount of shares: %d instead of %d", len(mnemonics), s.cfg.ShamirShares)
		s.logger.Error("msg", msg)
		err = errors.New(msg)
	}
	return
}

func (s *ShamirCoordinatorService) RecoverSeed(mnemonics []string) (seed string, err error) {
	selectedSharesLen := len(mnemonics)
	selectedSharesWords := make([][]uint16, selectedSharesLen)
	for index := 0; index < selectedSharesLen; index++ {
		selectedShareString := mnemonics[index]
		wordsInEachShare := len(strings.Fields(selectedShareString))
		resultWords, err := s.slip39Interface.WordsForStrings(selectedShareString, wordsInEachShare)
		if err != nil {
			return "", err
		}
		selectedSharesWords[index] = resultWords
	}
	secret, err := s.slip39Interface.Combine(selectedSharesWords, password)
	if err != nil {
		return
	}
	seed = hex.EncodeToString(secret)
	return
}
