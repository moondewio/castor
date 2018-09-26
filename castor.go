package castor

import (
	"fmt"
	"strconv"
)

// List lists all the PRs
func List(token string) error {
	prs, err := fetchPRs(token)
	if err != nil {
		return ExitErr(1, err)
	}

	for _, pr := range prs {
		printPR(pr)
	}

	return nil
}

// Review checksout the branch of a PR to review it, saving the status of the current
// branch to allow coming back to it later and continue with the work in progress.
func Review(n string, token string) error {
	prNum, err := strconv.Atoi(n)
	if err != nil {
		return ExitErrorF(1, "'%s' is not a number", n)
	}

	prs, err := fetchPRs(token)
	if err != nil {
		return ExitErr(1, err)
	}

	for _, pr := range prs {
		if pr.Number == prNum {
			return ExitErr(1, switchToPR(pr))
		}
	}

	return ExitErrorF(1, "PR #%v not found", prNum)
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
