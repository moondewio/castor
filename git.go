package castor

import (
	"fmt"
	"os"
	"strings"

	"github.com/whilp/git-urls"
)

var castorWIPMsg = "[CASTOR WIP]"
var castorWIPFile = ".castorwip"

func checkoutBranch(branch string) error {
	fmt.Printf("\nSwitching to branch `%s`\n\n", branch)
	if err := runWithPipe("git", "checkout", branch); err != nil {
		fmt.Println()
		if err := runWithPipe("git", "fetch"); err != nil {
			return err
		}
		fmt.Println()
		if err := runWithPipe("git", "checkout", branch); err != nil {
			return err
		}
	}
	fmt.Println()
	return nil
}

func switchToBranch(branch string) error {
	if !isRepo() {
		return fmt.Errorf("Not a git repository")
	}

	if isClean() {
		fmt.Print("Repository is clean, creating .castorwip to keep a reference to the branch\n\n")
		f, err := os.Create(castorWIPFile)
		if err != nil {
			return err
		}
		f.Close()
	}

	fmt.Printf("Saving Work In Progress\n\n")
	if err := runWithPipe("git", "stash", "save", "-u", castorWIPMsg); err != nil {
		fmt.Printf("\nCouldn't stash files...\n\n")
		return err
	}

	if err := checkoutBranch(branch); err != nil {
		fmt.Printf("\nFailed to checkout to branch `%s`, applying Work In Progress back\n\n", branch)
		if err := runWithPipe("git", "stash", "pop"); err != nil {
			fmt.Printf("\nFailed to apply changes...\n\n")
			return err
		}
		return err
	}

	if err := runWithPipe("git", "pull", "origin", branch); err != nil {
		fmt.Printf("\nSwitched to `%s` but failed to pull latest changes...\n", branch)
	} else {
		fmt.Printf("\nSwitched to branch `%s`\n", branch)
	}

	return nil
}

func goBack(branch string) error {
	if !isRepo() {
		return fmt.Errorf("Not a git repository")
	}

	cur, err := currentBranch()
	if err != nil {
		return err
	}

	if branch == cur {
		return fmt.Errorf("Already in branch `%s`", branch)
	}

	wip, ok := stashWIP(branch)
	if !ok {
		if branch == "" {
			return fmt.Errorf("Castor didn't save any Work In Progress in this repository")
		}
		return fmt.Errorf("Castor didn't save any Work In Progress in branch `%s`", branch)
	}

	if cur != wip.branch {
		fmt.Printf("Checkingout back to branch `%s`\n\n", wip.branch)

		err = runWithPipe("git", "checkout", wip.branch)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Recovering your Work In Progress\n\n")

	err = runWithPipe("git", "stash", "pop", wip.id)
	if err != nil {
		return err
	}

	if _, err := os.Stat(castorWIPFile); !os.IsNotExist(err) {
		err := os.Remove(castorWIPFile)
		if err != nil {
			return err
		}
		fmt.Print("\nRemoved .castorwip file\n\n")
		return runWithPipe("git", "status")
	}

	return nil
}

func currentBranch() (string, error) {
	return output("git", "rev-parse", "--abbrev-ref", "HEAD")
}

func isRepo() bool {
	return run("git", "rev-parse") == nil
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
	return output("git", "remote", "get-url", "origin")
}

func lastCommit() (string, error) {
	return output("git", "log", "--pretty=format:%s", "-n", "1")
}

type stashEntry struct {
	id     string
	branch string
	msg    string
}

// stashWIP finds the castor WIP branch by parcing the git stash output:
//
// $ git stash list
// stash@{0}: WIP on branch/current: a114cb6 Batman
// stash@{1}: On branch/current: [CASTOR WIP]
// stash@{2}: On branch/foo: b225dc7 foo
//
// If branch is an empty string, returns the last WIP branch.
func stashWIP(branch string) (stashEntry, bool) {
	stash, err := output("git", "stash", "list")
	if err != nil {
		return stashEntry{}, false
	}

	var match string
	for _, entry := range strings.Split(stash, "\n") {
		if strings.Contains(strings.TrimSpace(entry), castorWIPMsg) && strings.Contains(entry, branch) {
			match = entry
		}
	}

	if parts := strings.Split(match, ":"); match != "" && len(parts) >= 3 {
		return stashEntry{
			id:     strings.TrimSpace(parts[0]),
			branch: strings.Replace(strings.TrimSpace(parts[1]), "On ", "", 1),
			msg:    strings.TrimSpace(parts[2]),
		}, true
	}
	return stashEntry{}, false
}

func gitUser() (string, error) {
	return output("git", "config", "--global", "user.name")
}

func isClean() bool {
	out, _ := output("git", "status")

	return strings.Index(out, "nothing to commit") != -1
}
