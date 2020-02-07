package main

import (
	"context"
	"log"
	"net"
	"net/http/fcgi"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/kochurovro/redirectsrv/interface/repositories"
)

func main() {
	var cfg Config
	ctx := context.Background()
	c := memcache.New("localhost:11211")
	err := c.Ping()
	if err != nil {
		log.Fatal("main#memcached#init", err)
		return
	}

	a, err := c.Get("a")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(a)
	var wg sync.WaitGroup
	tasksCh := make(chan memcache.Item, 2)
	defer close(tasksCh)
	go pool(&wg, cfg, tasksCh, c)

	repo := repositories.NewSqlUrlRepo(ctx, "user:password@(localhost:3306)/test")
	r := gin.Default()
	r.GET("/TestApp", redirectHandler(cfg, repo, c, tasksCh))

	listener, _ := net.Listen("tcp", "127.0.0.1:5000")
	defer listener.Close()

	err = fcgi.Serve(listener, r)
	if err != nil {
		log.Fatal("main#", err)
	}
}
