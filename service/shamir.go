package service

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func (s *ShamirCoordinatorService) CreateMnemonics(hexSecret string) (mnemonics []string, err error) {
	// Define the command and arguments
	shamirScheme := strconv.Itoa(s.cfg.ShamirThreshold) + "of" + strconv.Itoa(s.cfg.ShamirShares)
	cmd := exec.Command(s.cfg.VirtualEnvPath+"/bin/python", s.cfg.VirtualEnvPath+"/bin/shamir", "create", "-S", hexSecret, shamirScheme)

	// Capture the output
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		fmt.Printf("stderr: %s\n", stderr.String())
		return
	}

	mnemonics = strings.Split(out.String(), "\n")

	// unwrap the result form the return message
	mnemonics = mnemonics[2:]
	mnemonics = mnemonics[:len(mnemonics)-1]

	if len(mnemonics) != s.cfg.ShamirShares {
		msg := fmt.Sprintf("The command didn't return the expected amount of shares: %d instead of %d", len(mnemonics), s.cfg.ShamirShares)
		fmt.Println(msg)
		err = errors.New(msg)
	}
	return
}

func (s *ShamirCoordinatorService) RecoverSeed(mnemonics []string) (seed string, err error) {
	args := []string{"../python/shamir_recover.py"}
	args = append(args, mnemonics...)

	cmd := exec.Command(s.cfg.VirtualEnvPath+"/bin/python", args...)
	// Capture the output
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		fmt.Printf("stderr: %s\n", stderr.String())
		fmt.Printf("stdout: %s\n", out.String())
		return
	}
	seed = out.String()
	seed = seed[:len(seed)-1]
	return
}
