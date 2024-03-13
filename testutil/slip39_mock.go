package testutil

import slip39 "github.com/rddl-network/bc-slip39-go"

type Slip39Mock struct {
}

func (s *Slip39Mock) Generate(_ uint8, _ []slip39.GroupDescriptor, _ []byte, _ string, _ uint8, _ *[0]byte) (result, wordsInEachShare int, sharesBuffer []uint16, err error) {
	return
}

func (s *Slip39Mock) Random() (random *[0]byte) {
	return
}
func (s *Slip39Mock) StringsForWords(_ []uint16, _ int) (result string, err error) {
	return
}

func (s *Slip39Mock) Combine(_ [][]uint16, _ string) (secret []byte, err error) {
	secret = []byte("00000000000000000000000000000000")
	return
}

func (s *Slip39Mock) WordsForStrings(_ string, _ int) (words []uint16, err error) {
	return
}
