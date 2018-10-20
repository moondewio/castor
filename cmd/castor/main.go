package main

import (
	"encoding/json"
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
		Name:  "prs",
		Usage: "List PRs",
		UsageText: strings.Join([]string{
			"$ castor prs --user other-user",
			"$ castor prs --closed --open=false",
			"$ castor prs --everyone",
			"$ castor prs --all",
		}, "\n   "),
		Action: func(ctx *cli.Context) error { return castor.List(loadConf(ctx)) },
		Flags:  prsFlags,
	},
	{
		Name:      "review",
		Usage:     "Checkout to a PR's branch to review it",
		UsageText: "$ castor review 42",
		Action:    reviewAction,
		Flags:     commonFlags,
	},
	{
		Name:      "back",
		Usage:     "Checkout to were you left off",
		UsageText: "$ castor back",
		Flags:     backFlags,
		Action:    func(ctx *cli.Context) error { return castor.GoBack(ctx.String("branch")) },
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
		Flags:  commonFlags,
	},
}

var tokenFlag = cli.StringFlag{
	Name:  "token",
	Usage: "GitHub API Token (repo and org:read permissions)",
}
var userFlag = cli.StringFlag{
	Name:  "user",
	Usage: "GitHub username",
}

var commonFlags = []cli.Flag{
	userFlag,
	tokenFlag,
}

var prsFlags = append(
	commonFlags,
	[]cli.Flag{
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
	}...)

var backFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "branch",
		Usage: "Branch to go back to",
	},
}

func reviewAction(ctx *cli.Context) error {
	args := ctx.Args()

	// TODO: prompt to input number (maybe list PRs?)
	if !args.Present() {
		return castor.ExitErrorF(1, "Missing PR number")
	}

	return castor.ReviewPR(ctx.Args().First(), loadConf(ctx))
}

func configAction(cxt *cli.Context) error {
	b, err := ioutil.ReadFile(castorfile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	conf := GlobalConf{}
	if len(b) > 0 {
		err = json.Unmarshal(b, &conf)
		if err != nil {
			return err
		}
	}

	lookUpFlags(&conf, cxt)

	b, err = json.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(castorfile, b, os.ModePerm)
}

// TODO: replace with spf13/viper
func loadConf(ctx *cli.Context) castor.Conf {
	c := config.NewConfig()
	err := c.Load(file.NewSource(file.WithPath(castorfile)))
	if err != nil {
		return castor.Conf{}
	}

	conf := GlobalConf{
		Token: c.Get("token").String(""),
		User:  c.Get("user").String(""),
	}
	lookUpFlags(&conf, ctx)

	if conf.User == "" {
		conf.User = castor.GitUser()
	}

	return castor.Conf{
		All:      ctx.Bool("all"),
		Everyone: ctx.Bool("everyone"),
		Closed:   ctx.Bool("closed"),
		Open:     ctx.Bool("open"),
		Token:    conf.Token,
		User:     conf.User,
	}
}

func lookUpFlags(conf *GlobalConf, ctx *cli.Context) {
	if ctx.String("token") != "" {
		conf.Token = ctx.String("token")
	}
	if ctx.String("user") != "" {
		conf.User = ctx.String("user")
	}
}
