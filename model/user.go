package model

import "github.com/google/go-github/github"

type UserProfile struct {
	user                github.User
	quarterCommitCount  map[string]uint
	langRepoCount       map[string]uint
	langStarCount       map[string]uint
	langCommitCount     map[string]uint
	repoCommitCount     map[string]uint
	repoStarCount       map[string]uint
	repoCommitCountDesc map[string]string
	repoStarCountDesc   map[string]string
	timeStamp           int64
}
