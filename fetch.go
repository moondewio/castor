package castor

import (
	"fmt"
	"net/http"

	"github.com/asmcos/requests"
)

func fetchPRs(token string) ([]PR, error) {
	owner, repo, err := ownerAndRepo()
	if err != nil {
		return []PR{}, err
	}

	r := requests.Requests()
	if token != "" {
		r.Header.Set("Authorization", "token "+token)
	}
	res, err := r.Get(githubPRURL(owner, repo))

	prs := []PR{}

	if err != nil {
		return prs, err
	}

	if res.R.StatusCode != http.StatusOK {
		return prs, fmt.Errorf("Failed to fetch, status: %v", res.R.StatusCode)
	}

	err = res.Json(&prs)

	return prs, err
}

func githubPRURL(owner, repo string) string {
	// GET /repos/:owner/:repo/pulls
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?status=open", owner, repo)
}
