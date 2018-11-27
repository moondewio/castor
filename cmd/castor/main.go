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

func init() {
	cur, err := user.Current()
	check(err)
	castorfile = path.Join(cur.HomeDir, ".castor.json")
}

func main() {
	app := cli.NewApp()

	app.Name = "castor"
	app.Version = "0.0.9"
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
			"This command requires a GitHub API Token to work.",
			"Check `castor help config` for more information.\n",
			"$ castor prs --user other-user",
			"$ castor prs --closed --open=false",
			"$ castor prs --everyone",
			"$ castor prs --all",
		}, "\n   "),
		Aliases: []string{"ls"},
		Action:  func(ctx *cli.Context) error { return castor.List(loadConf(ctx)) },
		Flags:   prsFlags,
	},
	{
		Name:  "review",
		Usage: "Checkout to a PR's branch to review it",
		UsageText: strings.Join([]string{
			"No need to save changes (i.e. stash or commit), castor takes care of that.\n",
			"IMPORTANT: castor uses `git stash` to save the Work In Progress, if you run",
			"`git stash drop` on it castor will not be able to go back to the branch",
			"and, most importatly, your Work In Progress will be lost.\n",
			"This command requires a GitHub API Token to work.",
			"Check `castor help config` for more information.\n",
			"$ castor review 42",
			"$ castor review 42 --no-stat",
		}, "\n   "),
		Aliases: []string{"r"},
		Action:  reviewAction,
		Flags:   reviewFlags,
	},
	{
		Name:  "back",
		Usage: "Go back to were you left off",
		UsageText: strings.Join([]string{
			"Goes back to the branch last brach `castor review x` was called from.",
			"Use `--branch [branch]` to go back to a particular branch.",
			"castor will recover the Work In Progress of the branch.\n",
			"$ castor back",
			"$ castor back --branch my-wip-branch",
		}, "\n   "),
		Aliases: []string{"b"},
		Flags:   backFlags,
		Action:  func(ctx *cli.Context) error { return castor.GoBack(ctx.String("branch")) },
	},
	{
		Name:  "config",
		Usage: "Save configuration to use with the other commands",
		UsageText: strings.Join([]string{
			"castor uses GitHub's APIv4, it requires a token, read more here:",
			"https://developer.github.com/v4/guides/forming-calls/#authenticating-with-graphql\n",
			"To search for private PRs and list users and teams involved in reviews",
			"castor requires 'repo' and 'org:read' permissions.",
			"If your PRs do not have team reviewers the 'org:read' permission can be skipped.",
			"If you are only searching public repos you and skip both permissions.",
			"Don't worry, the only thing castor does is search for PRs.\n",
			"$ castor config --token [token]",
			"$ castor config --user [github username]",
			"$ castor config --token [token] --user [github username]",
		}, "\n   "),
		Aliases: []string{"c"},
		Action:  configAction,
		Flags:   commonFlags,
	},
}

var tokenFlag = cli.StringFlag{
	Name:  "token",
	Usage: "GitHub API Token",
}
var userFlag = cli.StringFlag{
	Name:  "user",
	Usage: `GitHub username (default: "git config --global user.name")`,
}
var remoteFlag = cli.StringFlag{
	Name:  "remote",
	Value: "origin",
	Usage: "Repo remote",
}

var commonFlags = []cli.Flag{
	userFlag,
	tokenFlag,
}

var prsFlags = append(
	commonFlags,
	remoteFlag,
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
)

var reviewFlags = append(
	commonFlags,
	remoteFlag,
	cli.BoolFlag{
		Name:  "no-stat",
		Usage: "Don't show diff stats after changing branch",
	},
)

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

	conf := castor.Conf{}
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
	// make sure the castorfile is proerly created (i.e. `touch ~/.castor.json`)
	f, err := os.OpenFile(castorfile, os.O_RDONLY|os.O_CREATE, 0666)
	check(err)
	err = f.Close()
	check(err)

	c := config.NewConfig()
	err = c.Load(file.NewSource(file.WithPath(castorfile)))
	if err != nil {
		return castor.Conf{}
	}

	conf := castor.Conf{
		Token: c.Get("token").String(""),
		User:  c.Get("user").String(""),
	}
	lookUpFlags(&conf, ctx)
	flagsFallbacks(&conf)

	return conf
}

func lookUpFlags(conf *castor.Conf, ctx *cli.Context) {
	if ctx.String("token") != "" {
		conf.Token = ctx.String("token")
	}
	if ctx.String("user") != "" {
		conf.User = ctx.String("user")
	}

	conf.Remote = ctx.String("remote")

	conf.All = ctx.Bool("all")
	conf.Everyone = ctx.Bool("everyone")
	conf.Closed = ctx.Bool("closed")
	conf.Open = ctx.Bool("open")
	conf.ShowStats = !ctx.Bool("no-stat")
}

func flagsFallbacks(conf *castor.Conf) {
	if conf.User == "" {
		conf.User = castor.GitUser()
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
