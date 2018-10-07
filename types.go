package castor

import (
	"fmt"
	"time"
)

// ExitError implements the cli.ExitError interface
type ExitError struct {
	error
	code int
}

// ExitCode returns the error code
func (e ExitError) ExitCode() int {
	return e.code
}

// ExitErr returns an ExitError based on an error code an another error value.
func ExitErr(code int, err error) ExitError {
	switch v := err.(type) {
	// TODO: should we ignore `code` in this case?
	case ExitError:
		return v
	default:
		return ExitError{err, code}
	}
}

// ExitErrorF returns an ExitError based on an error code and a format specifier.
func ExitErrorF(code int, format string, a ...interface{}) ExitError {
	return ExitError{fmt.Errorf(format, a...), code}
}

// PR is the data of a PR provided by GitHub
type PR struct {
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	URL                string    `json:"url"`
	HTMLURL            string    `json:"html_url"`
	DiffURL            string    `json:"diff_url"`
	PatchURL           string    `json:"patch_url"`
	IssueURL           string    `json:"issue_url"`
	CommitsURL         string    `json:"commits_url"`
	ReviewCommentsURL  string    `json:"review_comments_url"`
	ReviewCommentURL   string    `json:"review_comment_url"`
	CommentsURL        string    `json:"comments_url"`
	StatusesURL        string    `json:"statuses_url"`
	Number             int       `json:"number"`
	State              string    `json:"state"`
	Title              string    `json:"title"`
	Body               string    `json:"body"`
	Locked             bool      `json:"locked"`
	ActiveLockReason   string    `json:"active_lock_reason"`
	Assignee           User      `json:"assignee"`
	Assignees          []User    `json:"assignees"`
	Labels             []Label   `json:"labels"`
	RequestedReviewers []User    `json:"requested_reviewers"`
	Milestone          Milestone `json:"milestone"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	ClosedAt           time.Time `json:"closed_at"`
	MergedAt           time.Time `json:"merged_at"`
	Head               Ref       `json:"head"`
	Base               Ref       `json:"base"`
	User               User      `json:"user"`
	Links              Links     `json:"_links"`
}

// Link contains the Href of a GitHub's API link
type Link struct {
	Href string `json:"href"`
}

// Links is a set of links included relevant to a GitHub's API response object
type Links struct {
	Self           Link `json:"self"`
	HTML           Link `json:"html"`
	Issue          Link `json:"issue"`
	Comments       Link `json:"comments"`
	ReviewComments Link `json:"review_comments"`
	ReviewComment  Link `json:"review_comment"`
	Commits        Link `json:"commits"`
	Statuses       Link `json:"statuses"`
}

// User contains the data of a user provided by GitHub's API
type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// Label is the GitHub's issue label
type Label struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

// Milestone contains GitHub's milestone data
type Milestone struct {
	URL          string    `json:"url"`
	HTMLURL      string    `json:"html_url"`
	LabelsURL    string    `json:"labels_url"`
	ID           int       `json:"id"`
	NodeID       string    `json:"node_id"`
	Number       int       `json:"number"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Creator      User      `json:"creator"`
	OpenIssues   int       `json:"open_issues"`
	ClosedIssues int       `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at"`
	DueOn        time.Time `json:"due_on"`
}

// Repo containes the data of a GitHub's repository
type Repo struct {
	ID               int         `json:"id"`
	NodeID           string      `json:"node_id"`
	Name             string      `json:"name"`
	FullName         string      `json:"full_name"`
	Owner            User        `json:"owner"`
	Private          bool        `json:"private"`
	HTMLURL          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ArchiveURL       string      `json:"archive_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BlobsURL         string      `json:"blobs_url"`
	BranchesURL      string      `json:"branches_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	CommentsURL      string      `json:"comments_url"`
	CommitsURL       string      `json:"commits_url"`
	CompareURL       string      `json:"compare_url"`
	ContentsURL      string      `json:"contents_url"`
	ContributorsURL  string      `json:"contributors_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	DownloadsURL     string      `json:"downloads_url"`
	EventsURL        string      `json:"events_url"`
	ForksURL         string      `json:"forks_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitURL           string      `json:"git_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	IssuesURL        string      `json:"issues_url"`
	KeysURL          string      `json:"keys_url"`
	LabelsURL        string      `json:"labels_url"`
	LanguagesURL     string      `json:"languages_url"`
	MergesURL        string      `json:"merges_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	PullsURL         string      `json:"pulls_url"`
	ReleasesURL      string      `json:"releases_url"`
	SSHURL           string      `json:"ssh_url"`
	StargazersURL    string      `json:"stargazers_url"`
	StatusesURL      string      `json:"statuses_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	TagsURL          string      `json:"tags_url"`
	TeamsURL         string      `json:"teams_url"`
	TreesURL         string      `json:"trees_url"`
	CloneURL         string      `json:"clone_url"`
	MirrorURL        string      `json:"mirror_url"`
	HooksURL         string      `json:"hooks_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Language         interface{} `json:"language"`
	ForksCount       int         `json:"forks_count"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Size             int         `json:"size"`
	DefaultBranch    string      `json:"default_branch"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	Topics           []string    `json:"topics"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	HasDownloads     bool        `json:"has_downloads"`
	Archived         bool        `json:"archived"`
	PushedAt         time.Time   `json:"pushed_at"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Permissions      struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
	AllowRebaseMerge bool `json:"allow_rebase_merge"`
	AllowSquashMerge bool `json:"allow_squash_merge"`
	AllowMergeCommit bool `json:"allow_merge_commit"`
	SubscribersCount int  `json:"subscribers_count"`
	NetworkCount     int  `json:"network_count"`
}

// Ref contains data of a GitHub's ref
type Ref struct {
	Label string `json:"label"`
	Ref   string `json:"ref"`
	Sha   string `json:"sha"`
	User  User   `json:"user"`
	Repo  Repo   `json:"repo"`
}

// -------------------------------------

type Search struct {
	IssueCount int   `json:"issueCount"`
	Nodes      []PR_ `json:"nodes"`
}

type PR_ struct {
	Number             int            `json:"number"`
	Title              string         `json:"title"`
	Author             WithLogin      `json:"author"`
	Closed             bool           `json:"closed"`
	BaseRefName        string         `json:"baseRefName"`
	HeadRefName        string         `json:"headRefName"`
	Labels             Labels         `json:"Labels"`
	SuggestedReviewers []Reviewer     `json:"suggestedReviewers"`
	ReviewRequests     ReviewRequests `json:"reviewRequests"`
	Reviews            Reviews        `json:"reviews"`
}

type WithLogin struct {
	Login string `json:"login"`
}

type Label_ struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Reviewer struct {
	Reviewer WithLogin `json:"reviewer"`
}

type RequestedReviewer struct {
	RequestedReviewer Foo `json:"requestedReviewer"`
}

type Foo struct {
	WithLogin
	Name string `json:"name"`
}

type Review_ struct {
	State       string    `json:"state"`
	SubmittedAt time.Time `json:"submittedAt"`
	URL         string    `json:"url"`
	Author      WithLogin `json:"Author"`
}

type ReviewRequests struct {
	TotalCount int                 `json:"totalCount"`
	Nodes      []RequestedReviewer `json:"Nodes"`
}

type Reviews struct {
	TotalCount int       `json:"totalCount"`
	Nodes      []Review_ `json:"Nodes"`
}

type Labels struct {
	TotalCount int      `json:"totalCount"`
	Nodes      []Label_ `json:"nodes"`
}
