package main

import (
	"fmt"
	"os/exec"
)

func switchToPR(pr PR) error {
	// TODO: improve logs (better feedback to the user)
	fmt.Printf("Switching to branch `%s`\n", pr.Head.Ref)
	fmt.Println("Saving work in progress ...")

	err := exec.Command("git", "add", ".").Run()
	if err != nil {
		return err
	}

	err = exec.Command("git", "commit", "-m", "'girp wip'").Run()
	if err != nil {
		fmt.Println("Failed to commit staged files, rolling back...")
		if rberr := exec.Command("git", "reset", ".").Run(); rberr != nil {
			fmt.Println("Failed to rollback staged files...")
			return rberr
		}
		return err
	}

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
	}

	return nil
}
