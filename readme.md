<h1 align="center">castor :hamster:</h1>

<p align="center">Review GitHub PRs and go back where you left of</p>

[![castor-v1.0.0](https://asciinema.org/a/205135.png)](https://asciinema.org/a/205135)

> **`castor` is still under development and I'm looking for more ways to improve
> the PR review process. All feedback and suggestions are welcome!!!**

## Install

```
$ go get -u github.com/moondewio/castor/cmd/castor
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
   $ castor config --token [token] --user [user]

VERSION:
   0.0.8

AUTHOR:
   Christian Gill (gillchristiang@gmail.com)

COMMANDS:
     prs, ls    List PRs
     review, r  Checkout to a PR's branch to review it
     back, b    Go back to were you left off
     config, c  Save configuration to use with the other commands
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
