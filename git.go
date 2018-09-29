package castor

import (
	"fmt"
	"os"
	"strings"

	"github.com/whilp/git-urls"
)

var castorWIPMsg = "[CASTOR WIP]"

func switchToPR(pr PR) error {
	// TODO: improve logs (better feedback to the user)
	err := setWipBranch()
	if err != nil {
		return err
	}

	fmt.Printf("Saving work in progress ...\n\n")

	err = runAndPipe("git", "stash", "save", "-u", castorWIPMsg)
	if err != nil {
		fmt.Printf("\nCouldn't stash files...\n\n")
		return err
	}

	fmt.Printf("\nSwitching to branch `%s`\n\n", pr.Head.Ref)

	// TODO: git fetch if pull fails (unless it's network?)
	err = runAndPipe("git", "checkout", pr.Head.Ref)
	if err != nil {
		fmt.Printf("\nFailed to checkout to branch `%s`, applying WIP changes\n\n", pr.Head.Ref)
		if rberr := runAndPipe("git", "stash", "pop"); rberr != nil {
			fmt.Printf("\nFailed to apply changes...\n\n")
			return rberr
		}
		return err
	}

	fmt.Println()

	err = runAndPipe("git", "pull", "origin", pr.Head.Ref)
	if err != nil {
		fmt.Printf("\nSwitched to `%s` but failed to pull latest changes...\n", pr.Head.Ref)
	} else {
		fmt.Printf("\nSwitched to `%s`...\n", pr.Head.Ref)
	}

	return nil
}

func goBack() error {
	wip, ok := stashWIP()
	if !ok {
		// TODO: improve this message
		return fmt.Errorf("No branch with Work In Progress.")
	}

	cur, err := currentBranch()
	if err != nil {
		return err
	}

	if cur != wip.branch {
		fmt.Printf("Checkingout back to branch `%s`\n\n", wip.branch)

		err = runAndPipe("git", "checkout", wip.branch)
		if err != nil {
			return err
		}
	}

	if wip.id != "" {

		fmt.Printf("Recovering your Work In Progress\n\n")

		return runAndPipe("git", "stash", "pop", wip.id)
	}

	return nil
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
