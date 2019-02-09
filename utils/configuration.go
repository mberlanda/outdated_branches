package utils

import "os"

// Config should be parsed from a JSON passed a argument later
type Config struct {
	OauthToken string `json:"oauth_token"`
	RepoAuthor string `json:"repo_author"`
	RepoName   string `json:"repo_name"`
}

func withDefault(a string, b string) string {
	if len(a) == 0 {
		return b
	}
	return a
}

func NewConfigFromEnv() Config {
	return Config{
		OauthToken: os.Getenv("GITHUB_OAUTH_TOKEN"),
		RepoAuthor: withDefault(os.Getenv("REPO_AUTHOR"), "mberlanda"),
		RepoName:   withDefault(os.Getenv("REPO_NAME"), "outdated_branches"),
	}
}
