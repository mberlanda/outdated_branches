package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Config should be parsed from a JSON passed a argument later
type Config struct {
	OauthToken string `json:"oauth_token"`
	RepoAuthor string `json:"repo_author"`
	RepoName   string `json:"repo_name"`
}

// https://developer.github.com/v3/pulls/#pull-requests
type GithubUser struct {
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

type GithubLabel struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

type GithubMilestone struct {
	URL          string     `json:"url"`
	HTMLURL      string     `json:"html_url"`
	LabelsURL    string     `json:"labels_url"`
	ID           int        `json:"id"`
	NodeID       string     `json:"node_id"`
	Number       int        `json:"number"`
	State        string     `json:"state"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Creator      GithubUser `json:"creator"`
	OpenIssues   int        `json:"open_issues"`
	ClosedIssues int        `json:"closed_issues"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ClosedAt     time.Time  `json:"closed_at"`
	DueOn        time.Time  `json:"due_on"`
}

type GithubTeam struct {
	ID              int         `json:"id"`
	NodeID          string      `json:"node_id"`
	URL             string      `json:"url"`
	Name            string      `json:"name"`
	Slug            string      `json:"slug"`
	Description     string      `json:"description"`
	Privacy         string      `json:"privacy"`
	Permission      string      `json:"permission"`
	MembersURL      string      `json:"members_url"`
	RepositoriesURL string      `json:"repositories_url"`
	Parent          interface{} `json:"parent"`
}

type GithubPermissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type GithubRepo struct {
	ID               int               `json:"id"`
	NodeID           string            `json:"node_id"`
	Name             string            `json:"name"`
	FullName         string            `json:"full_name"`
	Owner            GithubUser        `json:"owner"`
	Private          bool              `json:"private"`
	HTMLURL          string            `json:"html_url"`
	Description      string            `json:"description"`
	Fork             bool              `json:"fork"`
	URL              string            `json:"url"`
	ArchiveURL       string            `json:"archive_url"`
	AssigneesURL     string            `json:"assignees_url"`
	BlobsURL         string            `json:"blobs_url"`
	BranchesURL      string            `json:"branches_url"`
	CollaboratorsURL string            `json:"collaborators_url"`
	CommentsURL      string            `json:"comments_url"`
	CommitsURL       string            `json:"commits_url"`
	CompareURL       string            `json:"compare_url"`
	ContentsURL      string            `json:"contents_url"`
	ContributorsURL  string            `json:"contributors_url"`
	DeploymentsURL   string            `json:"deployments_url"`
	DownloadsURL     string            `json:"downloads_url"`
	EventsURL        string            `json:"events_url"`
	ForksURL         string            `json:"forks_url"`
	GitCommitsURL    string            `json:"git_commits_url"`
	GitRefsURL       string            `json:"git_refs_url"`
	GitTagsURL       string            `json:"git_tags_url"`
	GitURL           string            `json:"git_url"`
	IssueCommentURL  string            `json:"issue_comment_url"`
	IssueEventsURL   string            `json:"issue_events_url"`
	IssuesURL        string            `json:"issues_url"`
	KeysURL          string            `json:"keys_url"`
	LabelsURL        string            `json:"labels_url"`
	LanguagesURL     string            `json:"languages_url"`
	MergesURL        string            `json:"merges_url"`
	MilestonesURL    string            `json:"milestones_url"`
	NotificationsURL string            `json:"notifications_url"`
	PullsURL         string            `json:"pulls_url"`
	ReleasesURL      string            `json:"releases_url"`
	SSHURL           string            `json:"ssh_url"`
	StargazersURL    string            `json:"stargazers_url"`
	StatusesURL      string            `json:"statuses_url"`
	SubscribersURL   string            `json:"subscribers_url"`
	SubscriptionURL  string            `json:"subscription_url"`
	TagsURL          string            `json:"tags_url"`
	TeamsURL         string            `json:"teams_url"`
	TreesURL         string            `json:"trees_url"`
	CloneURL         string            `json:"clone_url"`
	MirrorURL        string            `json:"mirror_url"`
	HooksURL         string            `json:"hooks_url"`
	SvnURL           string            `json:"svn_url"`
	Homepage         string            `json:"homepage"`
	Language         interface{}       `json:"language"`
	ForksCount       int               `json:"forks_count"`
	StargazersCount  int               `json:"stargazers_count"`
	WatchersCount    int               `json:"watchers_count"`
	Size             int               `json:"size"`
	DefaultBranch    string            `json:"default_branch"`
	OpenIssuesCount  int               `json:"open_issues_count"`
	Topics           []string          `json:"topics"`
	HasIssues        bool              `json:"has_issues"`
	HasProjects      bool              `json:"has_projects"`
	HasWiki          bool              `json:"has_wiki"`
	HasPages         bool              `json:"has_pages"`
	HasDownloads     bool              `json:"has_downloads"`
	Archived         bool              `json:"archived"`
	PushedAt         time.Time         `json:"pushed_at"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Permissions      GithubPermissions `json:"permissions"`
	AllowRebaseMerge bool              `json:"allow_rebase_merge"`
	AllowSquashMerge bool              `json:"allow_squash_merge"`
	AllowMergeCommit bool              `json:"allow_merge_commit"`
	SubscribersCount int               `json:"subscribers_count"`
	NetworkCount     int               `json:"network_count"`
}

type GithubLink struct {
	Href string
}

type GithubLinks struct {
	Self           GithubLink `json:"self"`
	HTML           GithubLink `json:"html"`
	Issue          GithubLink `json:"issue"`
	Comments       GithubLink `json:"comments"`
	ReviewComments GithubLink `json:"review_comments"`
	ReviewComment  GithubLink `json:"review_comment"`
	Commits        GithubLink `json:"commits"`
	Statuses       GithubLink `json:"statuses"`
}

type GithubPullRequest struct {
	URL                string          `json:"url"`
	ID                 int             `json:"id"`
	NodeID             string          `json:"node_id"`
	HTMLURL            string          `json:"html_url"`
	DiffURL            string          `json:"diff_url"`
	PatchURL           string          `json:"patch_url"`
	IssueURL           string          `json:"issue_url"`
	CommitsURL         string          `json:"commits_url"`
	ReviewCommentsURL  string          `json:"review_comments_url"`
	ReviewCommentURL   string          `json:"review_comment_url"`
	CommentsURL        string          `json:"comments_url"`
	StatusesURL        string          `json:"statuses_url"`
	Number             int             `json:"number"`
	State              string          `json:"state"`
	Locked             bool            `json:"locked"`
	Title              string          `json:"title"`
	User               GithubUser      `json:"user"`
	Body               string          `json:"body"`
	Labels             []GithubLabel   `json:"labels"`
	Milestone          GithubMilestone `json:"milestone"`
	ActiveLockReason   string          `json:"active_lock_reason"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	ClosedAt           time.Time       `json:"closed_at"`
	MergedAt           time.Time       `json:"merged_at"`
	MergeCommitSha     string          `json:"merge_commit_sha"`
	Assignee           GithubUser      `json:"assignee"`
	Assignees          []GithubUser    `json:"assignees"`
	RequestedReviewers []GithubUser    `json:"requested_reviewers"`
	RequestedTeams     []GithubTeam    `json:"requested_teams"`
	Head               struct {
		Label string     `json:"label"`
		Ref   string     `json:"ref"`
		Sha   string     `json:"sha"`
		User  GithubUser `json:"user"`
		Repo  GithubRepo `json:"repo"`
	} `json:"head"`
	Base struct {
		Label string     `json:"label"`
		Ref   string     `json:"ref"`
		Sha   string     `json:"sha"`
		User  GithubUser `json:"user"`
		Repo  GithubRepo `json:"repo"`
	} `json:"base"`
	Links             GithubLinks `json:"_links"`
	AuthorAssociation string      `json:"author_association"`
}

// Subset of branch response for parsing purposes
// https://developer.github.com/v3/repos/branches/#get-branch

type GithubCommitAuthor struct {
	Name  string `json:"name"`
	Date  string `json:"date"`
	Email string `json:"email"`
}

type GithubCommit struct {
	URL    string `json:"url"`
	Sha    string `json:"sha"`
	NodeID string `json:"node_id"`
	Commit struct {
		Author  GithubCommitAuthor `json:"author"`
		URL     string             `json:"url"`
		Message string             `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		Committer    GithubCommitAuthor `json:"committer"`
		Verification struct {
			Verified  bool        `json:"verified"`
			Reason    string      `json:"reason"`
			Signature interface{} `json:"signature"`
			Payload   interface{} `json:"payload"`
		} `json:"verification"`
	} `json:"commit"`
	Author    GithubUser `json:"author"`
	Committer GithubUser `json:"committer"`
}

type GithubBranch struct {
	Name   string       `json:"name"`
	Commit GithubCommit `json:"commit"`
}

// https://developer.github.com/v3/repos/commits/#compare-two-commits
type GithubFile struct {
	Sha         string `json:"sha"`
	Filename    string `json:"filename"`
	Status      string `json:"status"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Changes     int    `json:"changes"`
	BlobURL     string `json:"blob_url"`
	RawURL      string `json:"raw_url"`
	ContentsURL string `json:"contents_url"`
	Patch       string `json:"patch"`
}

type GithubCommitCompare struct {
	URL             string       `json:"url"`
	HTMLURL         string       `json:"html_url"`
	PermalinkURL    string       `json:"permalink_url"`
	DiffURL         string       `json:"diff_url"`
	PatchURL        string       `json:"patch_url"`
	BaseCommit      GithubCommit `json:"base_commit"`
	MergeBaseCommit GithubCommit `json:"merge_base_commit"`
	Status          string       `json:"status"`
	AheadBy         int          `json:"ahead_by"`
	BehindBy        int          `json:"behind_by"`
	TotalCommits    int          `json:"total_commits"`
	// Commits         []GithubCommit `json:"commits"`
	// Files           []GithubFile `json:"files"`
}

type pullRequestList []GithubPullRequest

func (xs pullRequestList) concat(ys pullRequestList) pullRequestList {
	for _, y := range ys {
		xs = append(xs, y)
	}
	return xs
}

type AppMutex struct {
	lock          sync.Mutex
	BaseBranchMap map[string]string
	Client        *http.Client
	Config        *Config
}

func (a *AppMutex) authorize(req *http.Request) *http.Request {
	req.Header.Add("Authorization", "token "+a.Config.OauthToken)
	return req
}

func (a *AppMutex) apiOpenPullRequests(page int) *http.Request {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=open&page=%s", a.Config.RepoAuthor, a.Config.RepoName, strconv.Itoa(page))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "apiOpenPullRequests: "))
	}
	return req
}

func (a *AppMutex) apiHeadBranch(branch string) *http.Request {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s", a.Config.RepoAuthor, a.Config.RepoName, branch)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "apiHeadBranch: "))
	}
	return req
}

func (a *AppMutex) apiCommitCompare(base string, merge string) *http.Request {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/compare/%s...%s", a.Config.RepoAuthor, a.Config.RepoName, base, merge)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "apiCommitCompare: "))
	}
	return req
}

func (a *AppMutex) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := a.Client.Do(a.authorize(req))
	return resp, err
}

func (a *AppMutex) retrievePullRequestsWithPagination(page int) pullRequestList {
	pullRequests := pullRequestList{}
	respPr, errPr := a.doRequest(a.apiOpenPullRequests(page + 1))
	if errPr != nil {
		log.Fatal(errors.Wrap(errPr, "retrievePullRequestsWithPagination: "))
	}
	json.NewDecoder(respPr.Body).Decode(&pullRequests)
	defer respPr.Body.Close()
	if len(pullRequests) == 0 {
		return pullRequests
	}
	return pullRequests.concat(a.retrievePullRequestsWithPagination(page + 1))
}

func (a *AppMutex) cachedLastCommit(branchName string) (string, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	commit, found := a.BaseBranchMap[branchName]
	return commit, found
}

func (a *AppMutex) requestLastCommit(branchName string) string {
	a.lock.Lock()
	defer a.lock.Unlock()

	resp, err := a.doRequest(a.apiHeadBranch(branchName))
	if err != nil {
		log.Fatal(errors.Wrap(err, "requestLastCommit: "))
	}
	branch := GithubBranch{}
	json.NewDecoder(resp.Body).Decode(&branch)
	commit := branch.Commit.Sha
	// log.Print(fmt.Sprintf("Branch: %s | Last commit: %s", branchName, commit))
	a.BaseBranchMap[branchName] = branch.Commit.Sha
	return commit
}

func (a *AppMutex) getLastCommit(branchName string) string {
	commit, found := a.cachedLastCommit(branchName)
	if !found {
		return a.requestLastCommit(branchName)
	}
	return commit
}

func withDefault(a string, b string) string {
	if len(a) == 0 {
		return b
	}
	return a
}

func main() {
	log.Print("Started")

	config := Config{
		OauthToken: os.Getenv("GITHUB_OAUTH_TOKEN"),
		RepoAuthor: withDefault(os.Getenv("REPO_AUTHOR"), "mberlanda"),
		RepoName:   withDefault(os.Getenv("REPO_NAME"), "outdated_branches"),
	}
	if config.OauthToken == "" {
		log.Fatal("You need to export GITHUB_OAUTH_TOKEN env variable")
	}

	app := AppMutex{
		BaseBranchMap: make(map[string]string),
		Client:        &http.Client{},
		Config:        &config,
	}

	log.Print("Master branch commit Sha: " + app.getLastCommit("master"))

	pullRequests := app.retrievePullRequestsWithPagination(0)

	log.Print(strconv.Itoa(len(pullRequests)) + " Open Pull requests")

	fmt.Println("PR ID | Branch | Base Branch | CommitDiff | Created At")
	fmt.Println("------|--------|-------------|------------|-----------")

	eg := errgroup.Group{}
	for _, pr := range pullRequests {
		prID := "#" + strconv.Itoa(pr.Number)
		prCreatedAt := pr.CreatedAt.Format(time.UnixDate)
		headRef := pr.Head.Ref
		headSha := app.getLastCommit(headRef)
		baseRef := pr.Base.Ref
		baseSha := pr.Base.Sha
		eg.Go(func() error {
			compareCommit := GithubCommitCompare{}
			respCompare, errCompare := app.doRequest(app.apiCommitCompare(baseSha, headSha))
			if errCompare == nil {
				json.NewDecoder(respCompare.Body).Decode(&compareCommit)
				fmt.Println(prID + " | " + headRef + " | " + baseRef + " | " + strconv.Itoa(compareCommit.TotalCommits) + " | " + prCreatedAt)
			}
			return errCompare
		})
	}
	if err := eg.Wait(); err == nil {
		log.Print("Successfully retrieved all PRs.")
	} else {
		log.Fatal(errors.Wrap(err, "Received error:"))
	}

	log.Print("Finished")
}
