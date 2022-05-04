package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	addr := os.Getenv("REDIS_ADDR")
	if addr != "" {
		log.Printf("Connecting to redis: %s", addr)
	}

	c := redis.NewClient(&redis.Options{
		Addr:         addr,
		ReadTimeout:  50 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	for i := range make([]struct{}, 10) {
		wg.Add(1)
		log.Printf("Starting worker %d", i)

		go func(ctx context.Context, wg *sync.WaitGroup, c *redis.Client, i int) {
			defer wg.Done()

			key := fmt.Sprintf("m000:%d", i)
			c.Set(key, "hello", 5*time.Second) // just to give it something to find

			for {
				select {
				case <-ctx.Done():
					log.Printf("caught ctx.Done in %d", i)
					return
				default:
					if err := c.Get(key).Err(); err != nil {
						log.Printf("caught error in routine %d: %v", i, err)
					}
				}
			}
		}(ctx, &wg, c, i)
	}

	wg.Wait()
}
