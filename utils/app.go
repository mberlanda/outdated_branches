package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

type PullRequestList []GithubPullRequest

func (xs PullRequestList) concat(ys PullRequestList) PullRequestList {
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

func (a *AppMutex) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := a.Client.Do(a.authorize(req))
	return resp, err
}

func (a *AppMutex) ApiOpenPullRequests(page int) *http.Request {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=open&page=%s", a.Config.RepoAuthor, a.Config.RepoName, strconv.Itoa(page))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "ApiOpenPullRequests: "))
	}
	return req
}

func (a *AppMutex) ApiHeadBranch(branch string) *http.Request {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s", a.Config.RepoAuthor, a.Config.RepoName, branch)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "apiHeadBranch: "))
	}
	return req
}

func (a *AppMutex) ApiCommitCompare(base string, merge string) *http.Request {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/compare/%s...%s", a.Config.RepoAuthor, a.Config.RepoName, base, merge)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "apiCommitCompare: "))
	}
	return req
}

func (a *AppMutex) RetrievePullRequestsWithPagination(page int) PullRequestList {
	pullRequests := PullRequestList{}
	respPr, errPr := a.doRequest(a.ApiOpenPullRequests(page + 1))
	if errPr != nil {
		log.Fatal(errors.Wrap(errPr, "retrievePullRequestsWithPagination: "))
	}
	json.NewDecoder(respPr.Body).Decode(&pullRequests)
	defer respPr.Body.Close()
	if len(pullRequests) == 0 {
		return pullRequests
	}
	return pullRequests.concat(a.RetrievePullRequestsWithPagination(page + 1))
}

func (a *AppMutex) cachedLastCommit(branchName string) (string, bool) {
	a.lock.Lock()
	defer a.lock.Unlock()
	commit, found := a.BaseBranchMap[branchName]
	return commit, found
}

func (a *AppMutex) RequestLastCommit(branchName string) string {
	a.lock.Lock()
	defer a.lock.Unlock()

	resp, err := a.doRequest(a.ApiHeadBranch(branchName))
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

func (a *AppMutex) GetLastCommit(branchName string) string {
	commit, found := a.cachedLastCommit(branchName)
	if !found {
		return a.RequestLastCommit(branchName)
	}
	return commit
}

func (a *AppMutex) CompareCommits(baseSha string, headSha string) (*GithubCommitCompare, error) {
	compareCommit := GithubCommitCompare{}
	respCompare, errCompare := a.doRequest(a.ApiCommitCompare(baseSha, headSha))
	if errCompare != nil {
		return nil, errCompare
	}
	err := json.NewDecoder(respCompare.Body).Decode(&compareCommit)
	if err != nil {
		return nil, err
	}
	return &compareCommit, nil
}

func MakeAppWithDefaults() AppMutex {
	return AppMutex{
		BaseBranchMap: make(map[string]string),
		Client:        &http.Client{},
	}
}
