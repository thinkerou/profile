package server

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/hashicorp/golang-lru"
	"github.com/thinkerou/profile/model"
	"golang.org/x/oauth2"
)

var once sync.Once
var lruCache *lru.Cache

func init() {
	once.Do(func() {
		lruCache, _ = lru.New(1024)
	})
}

func GetUserProfile(c *gin.Context) {
	username := c.Param("user")
	data, ok := lruCache.Get(username)
	if ok {
		c.JSON(http.StatusOK, gin.H{"msg": data.(model.UserProfile)})
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

	var repoCommits = make(map[*github.Repository][]*github.RepositoryCommit)
	var langRepoGrouping = make(map[string][]*github.Repository)
	var repoStarCount = make(map[string]uint)
	var repoStarCountDesc = make(map[string]string)
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
		if *repo.WatchersCount > 0 {
			repoStarCount[*repo.Name] = uint(*repo.WatchersCount)
			repoStarCountDesc[*repo.Name] = *repo.Description
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

	var langRepoCount = make(map[string]uint)
	var langStarCount = make(map[string]uint)
	var langCommitCount = make(map[string]uint)
	for k, v := range langRepoGrouping {
		langRepoCount[k] = uint(len(v))
		for _, e := range v {
			langStarCount[k] += uint(*e.WatchersCount)
			langCommitCount[k] += uint(len(repoCommits[e]))
		}
	}

	var repoCommitCountDesc = make(map[string]string)
	var repoCommitCount = make(map[string]uint)
	for k, v := range repoCommits {
		repoCommitCount[*k.Name] = uint(len(v))
		repoCommitCountDesc[*k.Name] = *k.Description
	}

	uf := model.UserProfile{
		User:                *user,
		QuarterCommitCount:  nil,
		LangRepoCount:       langRepoCount,
		LangStarCount:       langStarCount,
		LangCommitCount:     langCommitCount,
		RepoCommitCount:     repoCommitCount,
		RepoStarCount:       repoStarCount,
		RepoCommitCountDesc: repoCommitCountDesc,
		RepoStarCountDesc:   repoStarCountDesc,
		TimeStamp:           time.Now().Unix(),
	}

	lruCache.Add(username, uf)
	c.JSON(http.StatusOK, gin.H{"msg": uf})
}
