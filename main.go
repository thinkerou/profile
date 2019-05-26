package main

import (
	"context"
	"fmt"
	// "net/http"
	"os"

	// "github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	/*r := gin.Default()
	r.GET("/welcome", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "Hello world!"})
	})
	r.Run()*/
	token := ""
	token = os.Getenv("GITHUB_TOKEN")
	println(token)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	orgs, _, err := client.Organizations.List(context.Background(), "thinkerou", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, org := range orgs {
		fmt.Printf("%v. %v\n", i+1, org.GetLogin())
	}
	/*repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	for i, repo := range repos {
		fmt.Printf("%v. %v\n", i+1, repo)
	}*/
}
