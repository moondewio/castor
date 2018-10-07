package castor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/asmcos/requests"
	"github.com/machinebox/graphql"
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
	res, err := r.Get(githubPRsURL(owner, repo))

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

// create a client (safe to share across requests)
var client = graphql.NewClient("https://api.github.com/graphql")

var prBranchNameQuery = `
query repoBranchName($owner: String!, $name:String!, $pr:Int!) {
  repository(owner: $owner, name: $name) {
    pullRequest(number:$pr) {
      headRefName
    }
  }
}
`

func getPRHeadName(id int, token string) (string, error) {
	owner, repo, err := ownerAndRepo()
	if err != nil {
		return "", err
	}

	req := graphql.NewRequest(prBranchNameQuery)
	req.Var("owner", owner)
	req.Var("name", repo)
	req.Var("pr", id)

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	var res map[string]map[string]map[string]string

	ctx := context.Background()

	if err := client.Run(ctx, req, &res); err != nil {
		return "", err
	}

	return res["repository"]["pullRequest"]["headRefName"], nil
}

var searchQuery = `
query search($query: String!) {
  search(query: $query, type: ISSUE, first: 10) {
    issueCount
    nodes {
      ... on PullRequest {
        number
        title
        author {
          login
        }
        closed
        baseRefName
        headRefName
        labels(first: 20) {
          totalCount
          nodes {
            name
            color
          }
        }
        suggestedReviewers {
          reviewer {
            login
          }
        }
        reviewRequests(first: 20) {
          totalCount
          nodes {
            requestedReviewer {
              ... on User {
                login
              }
              ... on Team {
                name
              }
            }
          }
        }
        reviews(first: 20) {
          totalCount
          nodes {
            state
            author {
              login
            }
            submittedAt
            url
          }
        }
      }
    }
  }
}
`

func searchPRs(user, token string) (Search, error) {
	owner, repo, err := ownerAndRepo()
	if err != nil {
		return Search{}, err
	}

	// make a request
	req := graphql.NewRequest(searchQuery)

	search := strings.Join([]string{
		"repo:" + owner + "/" + repo,
		"involves:" + user,
		"type:pr",
		"is:open",
		"is:unmerged",
	}, " ")

	// set any variables
	req.Var("query", search)

	// set header fields
	// req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "token "+token)

	// run it and capture the response
	var res struct {
		Search Search `json:"search"`
	}
	ctx := context.Background()

	if err := client.Run(ctx, req, &res); err != nil {
		return Search{}, err
	}

	return res.Search, nil
}

func githubPRsURL(owner, repo string) string {
	// GET /repos/:owner/:repo/pulls
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?status=open", owner, repo)
}

func githubPRURL(id int, owner, repo string) string {
	// GET /repos/:owner/:repo/pulls/:id
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%v", owner, repo, id)
}
