package castor

import (
	"fmt"
	"os"
	"os/exec"
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

	fmt.Println("Saving work in progress ...")

	err = exec.Command("git", "add", ".").Run()
	if err != nil {
		return err
	}

	err = exec.Command("git", "commit", "-m", castorWIPCommitMsg).Run()
	if err != nil {
		fmt.Println("Failed to commit staged files, rolling back...")
		if rberr := exec.Command("git", "reset", ".").Run(); rberr != nil {
			fmt.Println("Failed to rollback staged files...")
			return rberr
		}
		return err
	}

	fmt.Printf("Switching to branch `%s`\n", pr.Head.Ref)

	err = exec.Command("git", "checkout", pr.Head.Ref).Run()
	if err != nil {
		fmt.Printf("Failed to checkout to branch `%s`, reverting back\n", pr.Head.Ref)
		if rberr := exec.Command("git", "reset", "HEAD~").Run(); rberr != nil {
			fmt.Println("Failed to rollback commited files...")
			return rberr
		}
		return err
	}

	err = exec.Command("git", "pull", "origin", pr.Head.Ref).Run()
	if err != nil {
		fmt.Println("Success!!!")
		fmt.Printf("Switched to `%s` but failed pull lates changes...\n", pr.Head.Ref)
	} else {
		fmt.Println("Success!!!")
		fmt.Printf("Switched to `%s`...\n", pr.Head.Ref)
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

	fmt.Printf("Checkingout back to branch `%s`\n", wip)

	err = exec.Command("git", "checkout", wip).Run()
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

	fmt.Println("Recovering your Work In Progress")

	return exec.Command("git", "reset", "HEAD~").Run()
}

func currentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func isRepo() bool {
	return exec.Command("git", "rev-parse").Run() == nil
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
	out, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil {
		return "", err
	}

	if dir := strings.TrimSpace(string(out)); dir != ".git" {
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
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		// TODO: handle error properly, maybe dir has not .git/
		return "", err
	}
	return strings.Replace(string(output), "\n", "", 1), nil
}

func lastCommit() (string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%s", "-n", "1")
	output, err := cmd.Output()
	if err != nil {
		// TODO: handle error properly, maybe dir has not .git/
		return "", err
	}

	return string(output), nil
}
