package utils

import (
	"os/exec"
)

// CheckMockoonCli checks if mockoon-cli is available in the system path
func CheckMockoonCli() (bool, error) {
	cmd := exec.Command("mockoon-cli", "--version")
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
