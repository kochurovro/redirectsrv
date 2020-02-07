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
	ctx := context.Background()

	cfg, err := initConfig()
	if err != nil {
		log.Fatal("main#memcached#init", err)
		return
	}
	log.Println("cfg#", cfg)

	c := memcache.New(cfg.UrlMem)
	err = c.Ping()
	if err != nil {
		log.Fatal("main#memcached#init", err)
		return
	}

	var wg sync.WaitGroup
	tasksCh := make(chan memcache.Item, cfg.Workers)
	defer close(tasksCh)
	go pool(&wg, cfg, tasksCh, c)

	repo := repositories.NewSqlUrlRepo(ctx, cfg.UrlDB)
	r := gin.Default()
	r.GET("/TestApp", redirectHandler(cfg, repo, c, tasksCh))

	listener, _ := net.Listen("tcp", "127.0.0.1:5000")
	defer listener.Close()

	err = fcgi.Serve(listener, r)
	if err != nil {
		log.Fatal("main#", err)
	}
}
