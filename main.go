package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mberlanda/outdated_branches/utils"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	log.Print("Started")

	config := utils.NewConfigFromEnv()
	if config.OauthToken == "" {
		log.Fatal("You need to export GITHUB_OAUTH_TOKEN env variable")
	}

	app := utils.MakeAppWithDefaults()
	app.Config = &config

	log.Print("Master branch commit Sha: " + app.GetLastCommit("master"))

	pullRequests := app.RetrievePullRequestsWithPagination(0)

	log.Print(strconv.Itoa(len(pullRequests)) + " Open Pull requests")

	fmt.Println("PR ID | Branch | Base Branch | CommitDiff | Created At")
	fmt.Println("------|--------|-------------|------------|-----------")

	eg := errgroup.Group{}
	for _, pr := range pullRequests {
		prID := "#" + strconv.Itoa(pr.Number)
		prCreatedAt := pr.CreatedAt.Format(time.UnixDate)
		headRef := pr.Head.Ref
		headSha := app.GetLastCommit(headRef)
		baseRef := pr.Base.Ref
		baseSha := pr.Base.Sha
		eg.Go(func() error {
			compareCommit, errCompare := app.CompareCommits(baseSha, headSha)
			if errCompare == nil {
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
