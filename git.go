package castor

import (
	"fmt"
	"os"
	"strings"

	"github.com/whilp/git-urls"
)

var castorWIPCommitMsg = "[CASTOR WIP]"

func switchToPR(pr PR) error {
	// TODO: improve logs (better feedback to the user)
	err := setWipBranch()
	if err != nil {
		return err
	}

	fmt.Printf("Saving work in progress ...\n\n")

	err = runAndPipe("git", "add", ".")
	if err != nil {
		return err
	}

	err = runAndPipe("git", "commit", "-m", castorWIPCommitMsg)
	if err != nil {
		fmt.Printf("\nFailed to commit staged files, rolling back...\n\n")
		if rberr := runAndPipe("git", "reset", "."); rberr != nil {
			fmt.Printf("\nFailed to rollback staged files...\n\n")
			return rberr
		}
		return err
	}

	fmt.Printf("\nSwitching to branch `%s`\n\n", pr.Head.Ref)

	err = runAndPipe("git", "checkout", pr.Head.Ref)
	if err != nil {
		fmt.Printf("\nFailed to checkout to branch `%s`, reverting back\n\n", pr.Head.Ref)
		if rberr := runAndPipe("git", "reset", "HEAD~"); rberr != nil {
			fmt.Printf("\nFailed to rollback commited files...\n\n")
			return rberr
		}
		return err
	}

	fmt.Println()

	err = runAndPipe("git", "pull", "origin", pr.Head.Ref)
	if err != nil {
		fmt.Printf("\nSwitched to `%s` but failed to pull lates changes...\n", pr.Head.Ref)
	} else {
		fmt.Printf("\nSwitched to `%s`...\n", pr.Head.Ref)
	}

	return nil
}

// TODO: handle errors properly and display feedback
func goBack() error {
	wip, err := wipBranch()
	if err != nil {
		return err
	}

	cur, err := currentBranch()
	if err != nil {
		return err
	}

	if cur == wip {
		fmt.Printf("Already in branch `%s`\n", wip)
		return nil
	}

	fmt.Printf("Checkingout back to branch `%s`\n\n", wip)

	err = runAndPipe("git", "checkout", wip)
	if err != nil {
		return err
	}

	msg, err := lastCommit()
	if err != nil {
		return err
	}

	if msg != castorWIPCommitMsg {
		return nil
	}

	fmt.Printf("Recovering your Work In Progress\n\n")

	return runAndPipe("git", "reset", "HEAD~")
}

func currentBranch() (string, error) {
	return output("git", "rev-parse", "--abbrev-ref", "HEAD")
}

func isRepo() bool {
	return run("git", "rev-parse") == nil
}

// $ pwd
// /home/user/repo
//
// $ git rev-parse --git-dir
// .git

// $ pwd
// /home/user/repo/internal-dir
//
// $ git rev-parse --git-dir
// /home/user/repo/.git
func repoDir() (string, error) {
	out, err := output("git", "rev-parse", "--git-dir")
	if err != nil {
		return "", err
	}

	if dir := strings.TrimSpace(out); dir != ".git" {
		// TODO: should only replace /.git$/
		return strings.Replace(dir, ".git", "", 1), nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir, nil
}

func ownerAndRepo() (string, string, error) {
	rawurl, err := remoteURL()
	if err != nil {
		return "", "", err
	}

	return ownerAndRepoFromRemote(rawurl)
}

func ownerAndRepoFromRemote(remote string) (string, string, error) {
	url, err := giturls.Parse(remote)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(strings.Replace(url.Path, ".git", "", 1), "/")

	// TODO: handle len != 2 case (could be many things)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Cannot parse owner and repo from git remote origin")
	}

	return parts[0], parts[1], nil
}

func remoteURL() (string, error) {
	output, err := output("git", "remote", "get-url", "origin")
	return strings.Replace(output, "\n", "", 1), err
}

func lastCommit() (string, error) {
	return output("git", "log", "--pretty=format:%s", "-n", "1")
}
