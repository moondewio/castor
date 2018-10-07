package castor

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aybabtme/rgbterm"
	"github.com/lucasb-eyer/go-colorful"
)

// List lists all the PRs
func List(token string) error {
	prs, err := fetchPRs(token)
	if err != nil {
		return ExitErr(1, err)
	}

	printPRsTable(prs)

	return nil
}

// Review checksout the branch of a PR to review it, saving the status of the current
// branch to allow coming back to it later and continue with the work in progress.
func Review(n string, token string) error {
	prNum, err := strconv.Atoi(n)
	if err != nil {
		return ExitErrorF(1, "'%s' is not a number", n)
	}

	branch, err := getPRHeadName(prNum, token)
	if err != nil {
		return ExitErr(1, err)
	}

	err = switchToBranch(branch)
	if err != nil {
		return ExitErr(1, err)
	}

	return nil
}

// GoBack checkouts back to the last WIP brach
func GoBack() error {
	err := goBack()

	if err != nil {
		return ExitErr(1, err)
	}

	return nil
}

// Involves checks the PRs related to a user
func Involves(token string) error {
	user, err := gitUser()
	if err != nil {
		return ExitErr(1, err)
	}

	prs, err := searchPRs(user, token)
	if err != nil {
		return ExitErr(1, err)
	}

	printInvolvesTable(user, prs)

	return nil
}

func printInvolvesTable(user string, search Search) {
	if search.IssueCount == 0 {
		return
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 2, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, " PR\t TITLE\t BRANCH\t AUTHOR\t REVIEWS\t LABELS")

	for _, pr := range search.Nodes {
		var reviews string
		if pr.ReviewRequests.TotalCount > 0 {
			rev := "reviews"
			if pr.ReviewRequests.TotalCount == 1 {
				rev = "review "
			}
			var reviewers string
			for i, r := range pr.ReviewRequests.Nodes {
				reviewer := r.RequestedReviewer.Login
				if reviewer == "" {
					reviewer = r.RequestedReviewer.Name
				}
				if i == 0 {
					reviewers += reviewer
				} else {
					reviewers += ", " + reviewer
				}
			}
			reviews = fmt.Sprintf("Missing %v %s (%s)", pr.ReviewRequests.TotalCount, rev, reviewers)
		}

		fmt.Fprintf(
			w,
			" %v\t %s\t %s\t %s\t %s\t %s\n",
			pr.Number,
			truncate(pr.Title, 30),
			pr.HeadRefName,
			pr.Author.Login,
			reviews,
			labels2(pr.Labels),
		)
	}

	w.Flush()
}

func labels2(ls Labels) string {
	tags := make([]string, ls.TotalCount)

	for i, l := range ls.Nodes {
		r, g, b := hex2rgb("#" + l.Color)

		tags[i] = rgbterm.FgString(l.Name, r, g, b)
	}

	return strings.Join(tags, " ")
}

func printPRsTable(prs []PR) {
	if len(prs) == 0 {
		return
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 2, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, " PR\t TITLE\t BRANCH\t AUTHOR\t REVIEWS\t LABELS")

	for _, pr := range prs {
		var reviews string
		if len(pr.RequestedReviewers) > 0 {
			reviews = fmt.Sprintf("Missing %v reviews", len(pr.RequestedReviewers))
		}

		fmt.Fprintf(
			w,
			" %v\t %s\t %s\t %s\t %s\t %s\n",
			pr.Number,
			truncate(pr.Title, 30),
			pr.Head.Ref,
			pr.User.Login,
			reviews,
			labels(pr.Labels),
		)
	}

	w.Flush()
}

func truncate(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func hex2rgb(hex string) (uint8, uint8, uint8) {
	c, err := colorful.Hex(hex)
	if err != nil {
		// TODO: use a better default
		return 0, 0, 0
	}

	return c.RGB255()
}

func labels(ls []Label) string {
	tags := make([]string, len(ls))

	for i, l := range ls {
		r, g, b := hex2rgb("#" + l.Color)

		tags[i] = rgbterm.FgString(l.Name, r, g, b)
	}

	return strings.Join(tags, " ")
}
