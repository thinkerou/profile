package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/hashicorp/golang-lru"
	"github.com/jinzhu/now"
	"github.com/thinkerou/profile/model"
	"golang.org/x/oauth2"
)

var once sync.Once
var lruCache *lru.Cache
var userMap sync.Map

const EXPIRE = 60 // 3600 // one hour

func init() {
	once.Do(func() {
		lruCache, _ = lru.New(1024)
	})
}

func GetUserProfile(c *gin.Context) {
	username := c.Param("user")
	data, ok := lruCache.Get(username)
	if ok {
		t, o := userMap.Load(username)
		fmt.Println(o)
		fmt.Println(t)
		if o {
			if time.Now().Unix() - t.(int64) < EXPIRE {
				fmt.Println("return from cache")
				c.JSON(http.StatusOK, gin.H{"msg": data.(model.UserProfile)})
				return
			} else {
				fmt.Println("return expire")
				userMap.Delete(username)
			}
		}
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
			if repo.Description != nil {
				repoStarCountDesc[*repo.Name] = *repo.Description
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
		if k.Description != nil {
			repoCommitCountDesc[*k.Name] = *k.Description
		}
	}

	uf := model.UserProfile{
		User:                *user,
		QuarterCommitCount:  getCommitsForQuarters(user, repoCommits),
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
	userMap.Store(username, time.Now().Unix())

	c.JSON(http.StatusOK, gin.H{"msg": uf})
}

func getCommitsForQuarters(user *github.User, repoCommits map[*github.Repository][]*github.RepositoryCommit) map[string]uint {
	var result = make(map[string]uint)
	var allKeys = getQuarterNameList(*user.CreatedAt)
	var keys []string
	for _, v := range repoCommits {
		for _, e := range v {
			t := *e.Commit.Committer.Date
			name := getQuarterName(now.New(t).BeginningOfQuarter())
			keys = append(keys, name)
			if d, exist := result[name]; exist {
				result[name] = d + 1
			} else {
				result[name] = 1
			}
		}
	}
	for _, key := range allKeys {
		found := false
		for _, k := range keys {
			if key == k {
				found = true
				break
			}
		}
		if !found {
			result[key] = 0
		}
	}
	return result
}

func getQuarterNameList(ts github.Timestamp) []string {
	var result []string
	q := now.New(ts.Time).BeginningOfQuarter()
	m := now.New(time.Now()).BeginningOfQuarter()
	switch q.Month() {
	case 1:
		result = append(result, strconv.Itoa(q.Year())+"-Q1")
		fallthrough
	case 4:
		result = append(result, strconv.Itoa(q.Year())+"-Q2")
		fallthrough
	case 7:
		result = append(result, strconv.Itoa(q.Year())+"-Q3")
		fallthrough
	case 10:
		result = append(result, strconv.Itoa(q.Year())+"-Q4")
	}
	for i := q.Year() + 1; i < m.Year(); i++ {
		result = append(result, strconv.Itoa(i)+"-Q1")
		result = append(result, strconv.Itoa(i)+"-Q2")
		result = append(result, strconv.Itoa(i)+"-Q3")
		result = append(result, strconv.Itoa(i)+"-Q4")
	}
	switch m.Month() {
	case 10:
		result = append(result, strconv.Itoa(m.Year())+"-Q4")
		fallthrough
	case 7:
		result = append(result, strconv.Itoa(m.Year())+"-Q3")
		fallthrough
	case 4:
		result = append(result, strconv.Itoa(m.Year())+"-Q2")
		fallthrough
	case 1:
		result = append(result, strconv.Itoa(m.Year())+"-Q1")
	}
	return result
}

func getQuarterName(t time.Time) string {
	y := strconv.Itoa(t.Year())
	name := ""
	switch t.Month() {
	case 1:
		name = y + "-Q1"
	case 4:
		name = y + "-Q2"
	case 7:
		name = y + "-Q3"
	case 10:
		name = y + "-Q4"
	default:
		name = ""
	}
	return name
}
