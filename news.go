package main

import (
	"net/url"
	"time"
)

type NewsEntry struct {
	Title  string
	Url    *url.URL
	Time   time.Time
	Source string
}
