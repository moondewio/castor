package castor

import (
	"context"
	"strings"

	"github.com/machinebox/graphql"
)

func listOpenPRs(token string) (PRsSearch, error) {
	owner, repo, err := ownerAndRepo()
	if err != nil {
		return PRsSearch{}, err
	}
	searchQuery := strings.Join([]string{
		"repo:" + owner + "/" + repo,
		"type:pr",
		"is:open",
		"is:unmerged",
	}, " ")

	return searchPRs(token, searchQuery)
}

func searchPRsInvolvingUser(user, token string) (PRsSearch, error) {
	owner, repo, err := ownerAndRepo()
	if err != nil {
		return PRsSearch{}, err
	}
	searchQuery := strings.Join([]string{
		"repo:" + owner + "/" + repo,
		"involves:" + user,
		"type:pr",
		"is:open",
		"is:unmerged",
	}, " ")

	return searchPRs(token, searchQuery)
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

var prNodes = `
nodes {
  ... on PullRequest {
	number
	title
	url
	author {
	  login
	}
	headRefName
	labels(first: 20) {
	  totalCount
	  nodes {
		name
		color
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
  }
}
`

var listPRsQuery = `
query search($query: String!) {
  search(query: $query, type: ISSUE, first: 100) {
    issueCount
    ` + prNodes + `
  }
}
`

func searchPRs(token, searchQuery string) (PRsSearch, error) {
	req := graphql.NewRequest(listPRsQuery)
	req.Var("query", searchQuery)
	req.Header.Set("Authorization", "token "+token)

	var res struct {
		Search PRsSearch `json:"search"`
	}
	ctx := context.Background()

	if err := client.Run(ctx, req, &res); err != nil {
		return PRsSearch{}, err
	}

	return res.Search, nil
}
