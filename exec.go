package castor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func runWithPipe(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	color.New(color.FgYellow).Printf("$ %s %s\n\n", command, strings.Join(args, " "))

	return cmd.Run()
}

func run(command string, args ...string) error {
	return exec.Command(command, args...).Run()
}

func output(command string, args ...string) (string, error) {
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
