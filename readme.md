<h1 align="center">castor :hamster:</h1>

<p align="center">Review GitHub PRs and go back where you left of</p>

> **This is still WIP** :tm:
>
> **All feedback is welcome**

[![asciicast](https://asciinema.org/a/203682.png)](https://asciinema.org/a/203682)

## Install

```
$ go get github.com/gillchristian/castor/cmd/castor
```

## Use

```
NAME:
   castor - Review PRs in the terminal

USAGE:
   $ castor prs
   $ castor review 14
   $ castor back
   $ castor token [token]

VERSION:
   0.0.1

AUTHOR:
   Christian Gill (gillchristiang@gmail.com)

COMMANDS:
     prs      List all PRs
     review   Checkout to a PR's branch to review it
     back     Checkout to were you left off
     token    Save the GitHub API token to use with other commands
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value  GitHub API Token for accessing private repos
   --help, -h     show help
   --version, -v  print the version
```
