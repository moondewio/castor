<h1 align="center">castor :hamster:</h1>

<p align="center">Review GitHub PRs and go back where you left of</p>

> **This is still WIP** :tm:
>
> **All feedback is welcome**

[![castor-v1.0.0](https://asciinema.org/a/205135.png)](https://asciinema.org/a/205135)

## Install

```
$ go get github.com/moondewio/castor/cmd/castor
```

## Use

```
$ castor
NAME:
   castor - Review PRs in the terminal

USAGE:
   $ castor prs
   $ castor review 14
   $ castor back
   $ castor token [token]

VERSION:
   1.0.0

AUTHOR:
   Christian Gill (gillchristiang@gmail.com)

COMMANDS:
     prs      List PRs
     review   Checkout to a PR's branch to review it
     back     Checkout to were you left off
     token    Save the GitHub API token to use with other commands
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value  GitHub API Token for accessing private repos
   --help, -h     show help
   --version, -v  print the version
```

## Todo

- [x] Handle `git` errors: show command and output
- [ ] Handle `git` errors: better output
- [ ] Handle `git` errors: use command exit code
- [ ] `back`: support multiple branches
- [ ] `prs`: improve output
- [ ] `prs`: show all my PRs (in different repos)
- [ ] `review`: don't stash if there are no changes
- [ ] `review`: list changed files (with stats)
- [ ] Add tests
- [ ] Support different remotes than `origin`
- [ ] Add support for GitLab
- [ ] Token: `--token` flag should be per command, not global
