package wifi

import (
	"fmt"
	"os/exec"
)

func performCommand(command string, args ...string) ([]byte, error) {
	fmt.Print("Performing command: " + command + " ")
	fmt.Println(args)
	cmd := exec.Command(command, args...)
	return cmd.CombinedOutput()
}

func printOutput(output []byte, err error) {
	fmt.Println(string(output))
	if err != nil {
		fmt.Println(err)
	}
}
