package main

import (
	"testing"
	"time"
)

func TestFetchHackerNews(t *testing.T) {
	entries, err := fetchHackernews()
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}

	for _, entry := range entries {
		if entry.Title == "" {
			t.Fatalf("expected entry title to be non-empty, got: %+v", entry.Title)
		}
		if time.Time.IsZero(entry.Time) {
			t.Fatalf("expected entry time to be non-zero, got: %+v", entry.Time)
		}
		if entry.Source != "hackernews" {
			t.Fatalf("expected entry source to be hackernews, got: %+v", entry.Source)
		}
	}
}
