package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/moondewio/castor"
	"github.com/urfave/cli"
)

var castorfile string

// GlobalConf contains the app configuration
type GlobalConf struct {
	Token string `json:"token,omitempty"`
	User  string `json:"user,omitempty"`
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
	app.Version = "0.0.4"
	app.Author = "Christian Gill (gillchristiang@gmail.com)"
	app.Usage = "Review PRs in the terminal"
	app.UsageText = strings.Join([]string{
		"$ castor prs",
		"$ castor review 14",
		"$ castor back",
		"$ castor config --token [token] --user [user]",
	}, "\n   ")

	app.Commands = commands

	app.Run(os.Args)
}

var commands = []cli.Command{
	{
		Name:      "prs",
		Usage:     "List PRs",
		UsageText: "$ castor prs",
		Action:    func(c *cli.Context) error { return castor.List(loadConf(c)) },
		Flags:     prsFlags,
	},
	{
		Name:      "review",
		Usage:     "Checkout to a PR's branch to review it",
		UsageText: "$ castor review 14",
		Action:    reviewAction,
	},
	{
		Name:      "back",
		Usage:     "Checkout to were you left off",
		UsageText: "$ castor back",
		Flags:     backFlags,
		Action:    func(c *cli.Context) error { return castor.GoBack(c.String("branch")) },
	},
	{
		Name:  "config",
		Usage: "Save configuration to use with the other commands",
		UsageText: strings.Join([]string{
			"$ castor config --token [token]",
			"$ castor config --user [github username]",
			"$ castor config --token [token] --user [github username]",
		}, "\n   "),

		Action: configAction,
		Flags:  configFlags,
	},
}

var configFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "token",
		Usage: "GitHub API Token (repo and org:read permissions)",
	},
	cli.StringFlag{
		Name:  "user",
		Usage: "GitHub username",
	},
}

var prsFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "all",
		Usage: "All the projects I contribute to",
	},
	cli.BoolFlag{
		Name:  "everyone",
		Usage: "Include everyone's PRs, not only mine",
	},
	cli.BoolFlag{
		Name:  "closed",
		Usage: "Include closed PRs",
	},
	// cli.BoolTFlag defaults to true
	cli.BoolTFlag{
		Name:  "open",
		Usage: "Include open PRs (defaults to true)",
	},
}

var backFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "branch",
		Usage: "Branch to go back to",
	},
}

func reviewAction(c *cli.Context) error {
	args := c.Args()

	// TODO: prompt to input number (maybe list PRs?)
	if !args.Present() {
		return castor.ExitErrorF(1, "Missing PR number")
	}

	return castor.ReviewPR(c.Args().First(), loadConf(c))
}

func lookUpFlags(conf *map[string]string, c *cli.Context, flags ...string) {
	for _, flag := range flags {
		if c.String(flag) != "" {
			(*conf)[flag] = c.String(flag)
		}
	}
}

// TODO: don't override
func configAction(c *cli.Context) error {
	b, err := ioutil.ReadFile(castorfile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	conf := make(map[string]string)
	if len(b) > 0 {
		err = json.Unmarshal(b, &conf)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	lookUpFlags(&conf, c, "token", "user")

	b, err = json.Marshal(conf)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return ioutil.WriteFile(castorfile, b, os.ModePerm)
}

// TODO: replace with spf13/viper
// TODO: supports loading from flags (not only file)
func loadConf(c *cli.Context) castor.Conf {
	conf := config.NewConfig()
	err := conf.Load(file.NewSource(file.WithPath(castorfile)))
	if err != nil {
		return castor.Conf{}
	}

	return castor.Conf{
		All:      c.Bool("all"),
		Everyone: c.Bool("everyone"),
		Closed:   c.Bool("closed"),
		Open:     c.Bool("open"),
		Token:    conf.Get("token").String("token"),
		User:     conf.Get("user").String("user"),
	}
}
