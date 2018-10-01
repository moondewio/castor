package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/gillchristian/castor"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/urfave/cli"
)

var token string
var castorfile string

// Conf contains the app configuration
type Conf struct {
	Token string `json:"token"`
}

func init() {
	cur, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	castorfile = path.Join(cur.HomeDir, ".castor.json")
}

func main() {
	app := cli.NewApp()

	app.Name = "castor"
	app.Version = "1.0.0"
	app.Author = "Christian Gill (gillchristiang@gmail.com)"
	app.Usage = "Review PRs in the terminal"
	app.UsageText = strings.Join([]string{
		"$ castor prs",
		"$ castor review 14",
		"$ castor back",
		"$ castor token [token]",
	}, "\n   ")

	app.Commands = commands
	app.Flags = flags

	app.Run(os.Args)
}

var commands = []cli.Command{
	{
		Name:      "prs",
		Usage:     "List all PRs",
		UsageText: "$ castor prs",
		Action:    func(c *cli.Context) error { return castor.List(loadConf().Token) },
	},
	{
		Name:      "review",
		Usage:     "Checkout to a PR's branch to review it",
		UsageText: "$ castor review 14",
		Action:    reviewAction,
	},
	{
		// TODO: handle case of multiple WIPs
		// $ castor back [branch]
		// $ castor back # prompts branches when multiple WIPs
		Name:      "back",
		Usage:     "Checkout to were you left off",
		UsageText: "$ castor back",
		Action:    func(c *cli.Context) error { return castor.GoBack() },
	},
	{
		Name:  "token",
		Usage: "Save the GitHub API token to use with other commands",
		UsageText: strings.Join([]string{
			"$ castor token [token]",
			"$ castor --token [token] token",
		}, "\n   "),

		Action: tokenAction,
	},
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:        "token",
		Usage:       "GitHub API Token for accessing private repos",
		Destination: &token,
	},
}

func reviewAction(c *cli.Context) error {
	args := c.Args()

	if !args.Present() {
		return castor.ExitErrorF(1, "Missing PR number")
	}

	return castor.Review(c.Args().First(), loadConf().Token)
}

func tokenAction(c *cli.Context) error {
	if token != "" {
		return saveConf(Conf{Token: token})
	}

	args := c.Args()
	if !args.Present() {
		return castor.ExitErrorF(1, "No token provided")
	}

	return saveConf(Conf{Token: args.First()})
}

// TODO: create go-micro source for urfave/cli flags
func loadConf() Conf {
	if token != "" {
		return Conf{Token: token}
	}

	c := config.NewConfig()
	err := c.Load(file.NewSource(file.WithPath(castorfile)))
	if err != nil {
		return Conf{}
	}

	return Conf{Token: c.Get("token").String("token")}
}

func saveConf(conf Conf) error {
	content := []byte(`{"token": "` + conf.Token + `"}`)

	return ioutil.WriteFile(castorfile, content, os.ModeAppend)
}
