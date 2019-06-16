package model

import "github.com/google/go-github/github"

type UserProfile struct {
	User                github.User       `json:"user"`
	QuarterCommitCount  map[string]uint   `json:"quarterCommitCount"`
	LangRepoCount       map[string]uint   `json:"langRepoCount"`
	LangStarCount       map[string]uint   `json:"langStarCount"`
	LangCommitCount     map[string]uint   `json:"langCommitCount"`
	RepoCommitCount     map[string]uint   `json:"repoCommitCount"`
	RepoStarCount       map[string]uint   `json:"repoStarCount"`
	RepoCommitCountDesc map[string]string `json:"repoCommitCountDesc"`
	RepoStarCountDesc   map[string]string `json:"repoStarCountDesc"`
	TimeStamp           int64             `json:"timeStamp"`
}
