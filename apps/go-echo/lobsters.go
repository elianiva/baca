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

type LobstersEntry struct {
	ShortIdUrl string    `json:"short_id_url"`
	Title      string    `json:"title"`
	Url        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
}

const LobsterUrl = "https://lobste.rs/hottest.json"

func fetchLobsters(cache *bigcache.BigCache) ([]NewsEntry, error) {
	cached, err := cache.Get("lobsters")
	if err == nil {
		var entries []NewsEntry
		if err := json.Unmarshal(cached, &entries); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal cached lobsters")
		}
		return entries, nil
	}
	if !errors.Is(err, bigcache.ErrEntryNotFound) {
		return nil, errors.Wrap(err, "failed to get lobsters from cache")
	}

	resp, err := http.Get(LobsterUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch lobsters")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}(resp.Body)

	var entries []LobstersEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, errors.Wrap(err, "failed to decode lobsters response")
	}

	newsEntries := mapLobstersToNewsEntry(entries)
	serializedEntries, err := json.Marshal(newsEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal lobsters")
	}

	err = cache.Set("lobsters", serializedEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cache lobsters")
	}

	return newsEntries, nil
}

func mapLobstersToNewsEntry(entries []LobstersEntry) []NewsEntry {
	var newsEntries []NewsEntry
	for _, entry := range entries {
		rawUrl := entry.Url
		if rawUrl != "" {
			rawUrl = entry.ShortIdUrl
		}
		parsedUrl, err := url.Parse(rawUrl)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse url")
			continue
		}
		newsEntries = append(newsEntries, NewsEntry{
			Title:  entry.Title,
			Url:    parsedUrl,
			Time:   entry.CreatedAt,
			Source: "lobsters",
		})
	}
	return newsEntries
}
