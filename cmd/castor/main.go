package main

import (
	// "fmt"
	"os"
	"strings"

	"github.com/gillchristian/castor"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "castor"
	app.Version = "0.0.1"
	app.Author = "Christian Gill (gillchristiang@gmail.com)"
	app.Usage = "Review PRs in the terminal"
	app.UsageText = strings.Join([]string{
		"$ castor prs",
		"$ castor review 14",
	}, "\n   ")

	app.Commands = commands

	app.Run(os.Args)
}

var commands = []cli.Command{
	{
		Name:      "prs",
		Usage:     "List all PRs",
		UsageText: "$ castor prs",
		Action:    func(c *cli.Context) error { return castor.List() },
	},
	{
		Name:      "review",
		Usage:     "Checkout to a PR's branch to review it",
		UsageText: "$ castor review 14",
		Action:    review,
	},
}

func review(c *cli.Context) error {
	args := c.Args()

	if !args.Present() {
		return castor.ExitErrorF(1, "Missing PR number")
	}

	return castor.Review(c.Args().First())
}
