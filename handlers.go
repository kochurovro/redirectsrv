package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/kochurovro/redirectsrv/domain"
	"github.com/kochurovro/redirectsrv/interface/repositories"
)

type Params struct {
	Username string  `form:"username" binding:"required"`
	Ran      float64 `form:"ran" binding:"required"`
	PageURL  string  `form:"pageURL" binding:"required"`
}

func redirectHandler(cfg Config, repoURL *repositories.SqlUrlRepo, cache *memcache.Client, tasksCh chan<- memcache.Item) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p Params
		err := c.Bind(&p)
		if err != nil {
			c.JSON(502, gin.H{
				"error": err.Error(),
				"msg":   `You should use format params like "username=Anil&ran=0.97584378943&pageURL=http://xyz.com"`,
			})
			log.Println("redirectHandler#Bind", err)
			return
		}

		// check cache
		i, err := cache.Get(p.Username)
		if err != nil {
			if err != memcache.ErrCacheMiss {
				log.Println("redirectHandler#memcache#Get", err)
				c.JSON(502, gin.H{
					"error": err.Error(),
					"msg":   "Cache error",
				})
			}
		}
		if i != nil {
			c.Redirect(http.StatusFound, string(i.Value))
		}

		redirectUrl, err := repoURL.Get(p.Username)
		if err != nil {
			log.Println("redirectHandler#DB#Get", err)
			switch err {
			case domain.ErrRecordNotFound:
				// pageURL should be correct url
				if !IsUrl(p.PageURL) {
					c.JSON(502, gin.H{
						"error": err.Error(),
						"msg":   `You should use the correct format email`,
					})
					return
				}
				c.Redirect(http.StatusFound, p.PageURL)
				return
			case domain.ErrCloseConnection:
				c.Status(502)
				return
			default:
				c.Status(http.StatusInternalServerError)
			}
		}

		c.Redirect(http.StatusFound, redirectUrl)
		task := memcache.Item{
			Key:        p.Username,
			Value:      []byte("a"),
			Expiration: 100000,
		}
		tasksCh <- task
	}
}

// IsUrl checks if a string is correct url
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
