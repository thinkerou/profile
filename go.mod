module github.com/thinkerou/profile

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/hashicorp/golang-lru v0.5.1
	github.com/jinzhu/now v1.0.1
	golang.org/x/oauth2 v0.0.0-00010101000000-000000000000
)

replace golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190523182746-aaccbc9213b0
