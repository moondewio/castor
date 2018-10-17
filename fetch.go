package castor

import (
	"context"
	"strings"

	"github.com/machinebox/graphql"
)

func fetchPRs(conf PRsConfig, user, token string) (PRsSearch, error) {
	owner, repo, err := ownerAndRepo()
	if err != nil {
		return PRsSearch{}, err
	}
	search := []string{"type:pr"}

	if !conf.All {
		search = append(search, "repo:"+owner+"/"+repo)
	}
	if conf.Closed && !conf.Open {
		search = append(search, "is:closed")
	}
	if conf.Open && !conf.Closed {
		search = append(search, "is:open")
	}
	if !conf.Everyone {
		// TODO: involves vs author
		search = append(search, "involves:"+user)
	}

	return searchPRs(token, strings.Join(search, " "))
}

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
