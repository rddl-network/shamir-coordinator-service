package service

import slip39 "github.com/rddl-network/bc-slip39-go"

type ISlip39 interface {
	Generate(groupThreshold uint8, groups []slip39.GroupDescriptor, secret []byte, password string, iterationExponent uint8, randomGenerator *[0]byte) (result, wordsInEachShare int, sharesBuffer []uint16, err error)
	Random() (random *[0]byte)
	StringsForWords(words []uint16, wordsInEachShare int) (result string, err error)
	Combine(mnemonics [][]uint16, passphrase string) (secret []byte, err error)
	WordsForStrings(strings string, wordsInEachShare int) (words []uint16, err error)
}

type Slip39Interface struct {
}

func (s *Slip39Interface) Generate(groupThreshold uint8, groups []slip39.GroupDescriptor, secret []byte, password string, iterationExponent uint8, randomGenerator *[0]byte) (result, wordsInEachShare int, sharesBuffer []uint16, err error) {
	return slip39.Generate(groupThreshold, groups, secret, password, iterationExponent, randomGenerator)
}

func (s *Slip39Interface) Random() (random *[0]byte) {
	return slip39.Random()
}
func (s *Slip39Interface) StringsForWords(words []uint16, wordsInEachShare int) (result string, err error) {
	return slip39.StringsForWords(words, wordsInEachShare)
}

func (s *Slip39Interface) Combine(mnemonics [][]uint16, passphrase string) (secret []byte, err error) {
	return slip39.Combine(mnemonics, passphrase)
}

func (s *Slip39Interface) WordsForStrings(strings string, wordsInEachShare int) (words []uint16, err error) {
	return slip39.WordsForStrings(strings, wordsInEachShare)
}
