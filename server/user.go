package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/thinkerou/profile/model"
	"golang.org/x/oauth2"
)

func GetUserProfile(c *gin.Context) {
	username := c.Param("user")
	db, e := leveldb.OpenFile("user-profile", nil)
	if e != nil {
		println("1")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": e.Error()})
		return
	}
	defer db.Close()

	data, e := db.Get([]byte(username), nil)
	if e == nil {
		println("2")
		println(data)
		c.JSON(http.StatusOK, gin.H{"msg": data})
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	println(user.String())
	// c.JSON(http.StatusOK, gin.H{"email": *user.Email})

	var repos []*github.Repository
	i := 0
	for {
		i += 1
		opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{Page: i}}
		repo, _, err := client.Repositories.List(ctx, username, opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		if len(repo) == 0 {
			break
		}
		for _, e := range repo {
			if !*e.Fork && *e.Size != 0 {
				repos = append(repos, e)
			}
		}
	}
	println(len(repos))

	var repoCommits = make(map[*github.Repository][]*github.RepositoryCommit)
	var langRepoGrouping = make(map[string][]*github.Repository)
	for _, repo := range repos {
		if repo.Language == nil {
			r, ok := langRepoGrouping["Unknown"]
			if !ok {
				langRepoGrouping["Unknown"] = []*github.Repository{repo}
			} else {
				langRepoGrouping["Unknown"] = append(r, repo)
			}
		} else {
			r, ok := langRepoGrouping[*repo.Language]
			if !ok {
				langRepoGrouping[*repo.Language] = []*github.Repository{repo}
			} else {
				langRepoGrouping[*repo.Language] = append(r, repo)
			}
		}

		var cs []*github.RepositoryCommit
		i := 0
		for {
			i += 1
			opt := &github.CommitsListOptions{Author: username, ListOptions: github.ListOptions{Page: i}}
			commits, _, _ := client.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, opt)
			if len(commits) == 0 {
				repoCommits[repo] = cs
				break
			}
			cs = append(cs, commits...)
		}
	}
	println(len(repoCommits))
	println(len(langRepoGrouping))

	var langRepoCount = make(map[string]uint)
	// var langStarCount = make(map[string]uint)
	// var langCommitCount = make(map[string]uint)
	for k, v := range langRepoGrouping {
		langRepoCount[k] = uint(len(v))
		//
	}

	// var repoCommitCount = make(map[string]uint)
	// var repoStarCount = make(map[string]uint)

	uf := model.UserProfile{
		User:                *user,
		QuarterCommitCount:  nil,
		LangRepoCount:       langRepoCount,
		LangStarCount:       nil,
		LangCommitCount:     nil,
		RepoCommitCount:     nil,
		RepoCommitCountDesc: nil,
		RepoStarCountDesc:   nil,
		TimeStamp:           time.Now().Unix(),
	}

	// db.Put([]byte(username), []byte(uf), nil)
	db.Put([]byte(username), []byte("todo: user-profile"), nil)
	c.JSON(http.StatusOK, gin.H{"msg": uf})
}
