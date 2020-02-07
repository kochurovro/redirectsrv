package main

import (
	"log"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
)

func pool(wg *sync.WaitGroup, cfg *Config, tasksCh <-chan memcache.Item, client *memcache.Client) {
	wg.Add(cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		go worker(tasksCh, wg, client)
	}
}

func worker(tasksCh <-chan memcache.Item, wg *sync.WaitGroup, client *memcache.Client) {
	defer wg.Done()
	for {
		data, ok := <-tasksCh
		if !ok {
			log.Println("worker#exit")
			return
		}

		err := client.Add(&data)
		if err != nil {
			log.Println("worker#", err)
			continue
		}
		err = client.Set(&data)
		if err != nil {
			log.Println("worker#", err)
			continue
		}
	}
}
