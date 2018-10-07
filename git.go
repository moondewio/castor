package castor

import (
	"fmt"
	"strings"

	"github.com/whilp/git-urls"
)

var castorWIPMsg = "[CASTOR WIP]"

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

func switchToPR(pr PR) error {
	if !isRepo() {
		return fmt.Errorf("Not a git repository")
	}

	fmt.Printf("Saving Work In Progress\n\n")
	if err := runWithPipe("git", "stash", "save", "-u", castorWIPMsg); err != nil {
		fmt.Printf("\nCouldn't stash files...\n\n")
		return err
	}

	if err := checkoutBranch(pr.Head.Ref); err != nil {
		fmt.Printf("\nFailed to checkout to branch `%s`, applying Work In Progress back\n\n", pr.Head.Ref)
		if err := runWithPipe("git", "stash", "pop"); err != nil {
			fmt.Printf("\nFailed to apply changes...\n\n")
			return err
		}
		return err
	}

	if err := runWithPipe("git", "pull", "origin", pr.Head.Ref); err != nil {
		fmt.Printf("\nSwitched to `%s` but failed to pull latest changes...\n", pr.Head.Ref)
	} else {
		fmt.Printf("\nSwitched to `%s`...\n", pr.Head.Ref)
	}

	return nil
}

func goBack() error {
	if !isRepo() {
		return fmt.Errorf("Not a git repository")
	}

	wip, ok := stashWIP()
	if !ok {
		return fmt.Errorf("Castor didn't save any Work In Progress in this repository")
	}

	cur, err := currentBranch()
	if err != nil {
		return err
	}

	if cur != wip.branch {
		fmt.Printf("Checkingout back to branch `%s`\n\n", wip.branch)

		err = runWithPipe("git", "checkout", wip.branch)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Recovering your Work In Progress\n\n")

	return runWithPipe("git", "stash", "pop", wip.id)
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

// $ stash list
// stash@{0}: WIP on branch/current: a114cb6 Batman
// stash@{1}: On branch/current: [CASTOR WIP]
// stash@{2}: On branch/foo: b225dc7 foo
func stashWIP() (stashEntry, bool) {
	stash, err := output("git", "stash", "list")
	if err != nil {
		return stashEntry{}, false
	}

	var match string
	for _, entry := range strings.Split(stash, "\n") {
		if strings.Contains(strings.TrimSpace(entry), castorWIPMsg) {
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
