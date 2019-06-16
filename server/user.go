package server

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetUserProfile(c *gin.Context) {
	username := c.Param("user")

	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	users, _, err := client.Users.Get(ctx, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	println(users.String())
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
}

