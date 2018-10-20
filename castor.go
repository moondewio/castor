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

// Conf holds the configuration for listing PRs.
type Conf struct {
	All, Everyone, Closed, Open bool
	Token, User                 string
}

// List lists PRs
func List(conf Conf) error {
	prs, err := fetchPRs(conf)
	if err != nil {
		return ExitErr(1, err)
	}

	printPRsList(prs.IssueCount, prs.Nodes, conf)

	return nil
}

// ReviewPR checksout the branch of a PR to review it, saving the status of the current
// branch to allow coming back to it later and continue with the work in progress.
func ReviewPR(n string, conf Conf) error {
	prNum, err := strconv.Atoi(n)
	if err != nil {
		return ExitErrorF(1, "'%s' is not a number", n)
	}

	branch, err := getPRHeadName(prNum, conf)
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
func GoBack(branch string) error {
	err := goBack(branch)

	if err != nil {
		return ExitErr(1, err)
	}

	return nil
}

// TODO: don't print status if all open (`--closed` could be merged/closed)
func printPRsList(count int, prs []SearchPR, conf Conf) {
	if count == 0 {
		return
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 2, 1, ' ', tabwriter.Debug)
	switch {
	case conf.All:
		fmt.Fprintln(w, " PR\t REPO\t TITLE\t BRANCH\t AUTHOR\t STATUS\t REVIEWS\t LABELS")
	default:
		fmt.Fprintln(w, " PR\t TITLE\t BRANCH\t AUTHOR\t STATUS\t REVIEWS\t LABELS")
	}

	for _, pr := range prs {
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

		// TODO: fix string len when using colors (breaks column width)
		status := "Open" // rgbterm.FgString("Open", 0, 255, 0)
		if pr.Closed {
			status = "Closed" // rgbterm.FgString("Closed", 255, 0, 0)
		}
		if pr.Merged {
			status = "Merged" // rgbterm.FgString("Merged", 111, 66, 193)
		}

		switch {
		case conf.All:
			fmt.Fprintf(
				w,
				" %v\t %s\t %s\t %s\t %s\t %s\t %s\t %s\n",
				pr.Number,
				pr.HeadRepositoryOwner.Login+"/"+pr.HeadRepository.Name,
				truncate(pr.Title, 30),
				truncate(pr.HeadRefName, 30),
				pr.Author.Login,
				status,
				reviews,
				labels(pr.Labels),
			)
		default:
			fmt.Fprintf(
				w,
				" %v\t %s\t %s\t %s\t %s\t %s\t %s\n",
				pr.Number,
				truncate(pr.Title, 30),
				truncate(pr.HeadRefName, 30),
				pr.Author.Login,
				status,
				reviews,
				labels(pr.Labels),
			)
		}
	}

	w.Flush()
}

func labels(ls Labels) string {
	tags := make([]string, ls.TotalCount)

	for i, l := range ls.Nodes {
		r, g, b := hex2rgb("#" + l.Color)

		tags[i] = rgbterm.FgString(l.Name, r, g, b)
	}

	return strings.Join(tags, " ")
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
