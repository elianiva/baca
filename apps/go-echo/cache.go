package main

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/rs/zerolog/log"
)

func newCache() *bigcache.BigCache {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         10 * time.Minute,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   8192,
	}

	contextTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cache, initErr := bigcache.New(contextTimeout, config)
	if initErr != nil {
		log.Fatal().Err(initErr).Msg("failed to initialise bigcache")
	}
	
	return cache
}
