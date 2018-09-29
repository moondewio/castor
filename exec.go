package castor

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func runWithPipe(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	c := color.New(color.FgYellow)
	c.Printf("$ %s %s\n\n", command, strings.Join(args, " "))
	err = cmd.Start()
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, stdout)
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stderr, stderr)
	return err
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
