package model

import "github.com/google/go-github/github"

type UserProfile struct {
	user                github.User       `json:"user"`
	quarterCommitCount  map[string]uint   `json:"quarterCommitCount"`
	langRepoCount       map[string]uint   `json:"langRepoCount"`
	langStarCount       map[string]uint   `json:"langStarCount"`
	langCommitCount     map[string]uint   `json:"langCommitCount"`
	repoCommitCount     map[string]uint   `json:"repoCommitCount"`
	repoStarCount       map[string]uint   `json:"repoStarCount"`
	repoCommitCountDesc map[string]string `json:"repoCommitCountDesc"`
	repoStarCountDesc   map[string]string `json:"repoStarCountDesc"`
	timeStamp           int64             `json:"timeStamp"`
}
