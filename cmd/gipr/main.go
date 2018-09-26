package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/whilp/git-urls"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "prs":
		listPRs()
	case "view":
		if len(os.Args) >= 2 {
			viewPR(os.Args[2])
		} else {
			fmt.Println("Missing PR number...")
		}
	}
}

func listPRs() {
	owner, repo, err := ownerAndRepo()
	if err != nil {
		fmt.Println("There was an error, exiting")
		log.Fatal(err)
	}

	prs, err := fetchPRs(owner, repo)
	if err != nil {
		fmt.Println("There was an error, exiting")
		log.Fatal(err)
	}

	for _, pr := range prs {
		printPR(pr)
	}
}

func viewPR(n string) {
	prNum, err := strconv.Atoi(n)
	if err != nil {
		fmt.Println("That's not a numner!")
		log.Fatal(err)
	}

	owner, repo, err := ownerAndRepo()
	if err != nil {
		fmt.Println("Could not parse remote (origin) URL...")
		log.Fatal(err)
	}

	prs, err := fetchPRs(owner, repo)
	if err != nil {
		fmt.Println("Failed to fetch PRs information...")
		log.Fatal(err)
	}

	for _, pr := range prs {
		if pr.Number == prNum {
			err := switchToPR(pr)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	fmt.Printf("PR #%v not found", prNum)
}

func printPR(pr PR) {
	fmt.Println("Title: ", pr.Title)
	fmt.Println("URL: ", pr.URL)
	fmt.Println("Assignee: ", pr.Assignee.Login)
	fmt.Println("User: ", pr.User.Login)
	fmt.Println("Branch:", pr.Head.Ref)
	fmt.Println("Base branch:", pr.Base.Ref)
	fmt.Println("PR #", pr.Number)
	fmt.Printf("Missing %v reviews\n", len(pr.RequestedReviewers))
	fmt.Println("---")
}

func ownerAndRepo() (string, string, error) {
	rawurl, err := remoteURL()
	if err != nil {
		return "", "", err
	}

	return ownerAndRepoFromRemote(rawurl)
}

func remoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Replace(string(output), "\n", "", 1), nil
}

func ownerAndRepoFromRemote(remote string) (string, string, error) {
	url, err := giturls.Parse(remote)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(strings.Replace(url.Path, ".git", "", 1), "/")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("Cannot get owner and repo from git remote origin")
	}

	return parts[0], parts[1], nil
}
