package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

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

func fetchHackernews() ([]NewsEntry, error) {
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

	return mapHackernewsToNewsEntry(entries), nil
}

func mapHackernewsToNewsEntry(entries []HackernewsEntry) []NewsEntry {
	var newsEntries []NewsEntry
	for _, entry := range entries {
		parsedUrl, err := url.Parse(entry.URL)
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
