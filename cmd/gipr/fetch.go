package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchPRs(owner, repo string) ([]PR, error) {
	// GET /repos/:owner/:repo/pulls
	res, err := http.Get(githubPRURL(owner, repo))

	prs := []PR{}

	if err != nil {
		return prs, err
	}
	if res.StatusCode != http.StatusOK {
		return prs, fmt.Errorf("Failed to fetch, status: %v", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return prs, err
	}

	err = json.Unmarshal(b, &prs)
	if err != nil {
		return prs, err
	}

	return prs, nil
}

func githubPRURL(owner, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?status=open", owner, repo)
}
