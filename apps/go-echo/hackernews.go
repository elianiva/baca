package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type HackernewsEntry struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Time  int64  `json:"time"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

const HnUrl = "https://api.hackerwebapp.com/news?page=1"

func fetchHackernews(cache *bigcache.BigCache) ([]NewsEntry, error) {
	cached, err := cache.Get("hackernews")
	if err == nil {
		var entries []NewsEntry
		if err := json.Unmarshal(cached, &entries); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal cached hackernews")
		}
		return entries, nil
	}
	if !errors.Is(err, bigcache.ErrEntryNotFound) {
		return nil, errors.Wrap(err, "failed to get hackernews from cache")
	}

	resp, err := http.Get(HnUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch hackernews")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}(resp.Body)

	var entries []HackernewsEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, errors.Wrap(err, "failed to decode hackernews response")
	}

	newsEntries := mapHackernewsToNewsEntry(entries)
	serializedEntries, err := json.Marshal(newsEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal hackernews")
	}

	err = cache.Set("hackernews", serializedEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cache hackernews")
	}

	return newsEntries, nil
}

func mapHackernewsToNewsEntry(entries []HackernewsEntry) []NewsEntry {
	var newsEntries []NewsEntry
	for _, entry := range entries {
		parsedUrl, err := url.Parse(entry.URL)
		if parsedUrl.Host == "" {
			parsedUrl.Scheme = "https"
			parsedUrl.Host = "news.ycombinator.com"
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to parse url")
			continue
		}

		newsEntries = append(newsEntries, NewsEntry{
			Title:  entry.Title,
			Url:    parsedUrl,
			Time:   time.Unix(entry.Time, 0).UTC(),
			Source: "hackernews",
		})
	}
	return newsEntries
}
