<h1 align="center">castor :hamster:</h1>

<p align="center">Review GitHub PRs and go back where you left of</p>

> **This is still WIP** :tm:
>
> **All feedback is welcome**

[![castor-v1.0.0](https://asciinema.org/a/205135.png)](https://asciinema.org/a/205135)

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
   1.0.0

AUTHOR:
   Christian Gill (gillchristiang@gmail.com)

COMMANDS:
     prs      List PRs
     review   Checkout to a PR's branch to review it
     back     Checkout to were you left off
     config   Save configuration to use with the other commands
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
