package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/septemhill/dh/shorturl"
	"github.com/septemhill/dh/shorturl/cache"
	"github.com/septemhill/dh/shorturl/endpoints"
	"github.com/septemhill/dh/shorturl/repository"
	"github.com/septemhill/dh/shorturl/service"
)

func main() {
	ctx := context.Background()
	repo := repository.NewRepository(os.Getenv(PostgresConn))
	cache := cache.NewCache(os.Getenv(RedisAddress), cache.Password(os.Getenv(RedisPassword)))

	op := os.Getenv(ShortURLServerPort)
	port, err := strconv.ParseInt(op, 10, 64)
	if err != nil {
		log.Fatal("failed to get short url server port: " + err.Error())
	}

	srv := shorturl.NewServer(int(port), endpoints.MakeEndpoints(service.NewShortURLService(ctx, repo, cache)))
	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
