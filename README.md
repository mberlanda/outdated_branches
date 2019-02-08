# Outdated Branches

When you are in release management duty, you may want to have an overview of the current rebase state of all open prs.
For this purpose, you can use outdated branches`.

**What would you need?**

* `go` installed on your machine (tested against 1.11.4)
* a github oauth token (`Settings > Developer settings > Personal access tokens` and may Enable SSO)

**Usage:**

```sh
$ GITHUB_OAUTH_TOKEN=my-token go run main.go # by default mberlanda/outdated_branches
$ GITHUB_OAUTH_TOKEN=my-token REPO_NAME=cheidelacoriera go run main.go # mberlanda/outdated_branches
$ GITHUB_OAUTH_TOKEN=my-token REPO_AUTHOR=rails REPO_NAME=rails go run main.go # rails/rails
```

**Output:**

When tested against rails/rails:

```
2019/01/29 18:17:09 Started
2019/01/29 18:17:09 Master branch commit Sha: 3d22069c6355dc60be65e01958cf32917bc53142
2019/01/29 18:17:48 729 Open Pull requests
Branch | BranchSha | CommitDiff
# ...
active-storage-add-proxying-and-direct-downloads | 00a38fc98858e1153004c4dfaf0dd8bf8d65bec3 | 0
deprecation-warning-for-store-attributes | c3b90004f3b9f55380e50a2eee48a61f54446a43 | 0
2019/01/29 18:18:00 Finished
```

**Next steps:**

* reading config from a json file
* choose output format and content (currenty printing to STDOUT branch name, base commit id, number of commits between base and master)
* eventually add automated actions (e.g. add a label to the PR)
