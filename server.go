package main

import (
	"net/http"

	"github.com/allegro/bigcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type ServerConfig struct {
	Cache *bigcache.BigCache
}

func newServer(config *ServerConfig) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewRenderer()

	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		source := c.QueryParam("source")
		if source == "" {
			source = "hackernews"
		}

		entries, err := fetchNews(config.Cache, source)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.Render(http.StatusOK, "index.gohtml", map[string]any{
			"Entries": entries,
			"Source":  source,
		})
	})

	e.GET("/source/:source", func(c echo.Context) error {
		source := c.Param("source")
		if source == "" {
			source = "hackernews"
		}

		entries, err := fetchNews(config.Cache, source)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.Render(http.StatusOK, "source-update.gohtml", map[string]any{
			"Entries": entries,
			"Source":  source,
		})
	})

	return e
}

func fetchNews(cache *bigcache.BigCache, source string) (entries []NewsEntry, err error) {
	switch source {
	case "lobsters":
		entries, err = fetchLobsters(cache)
		return
	case "hackernews":
		entries, err = fetchHackernews(cache)
		return
	default:
		err = errors.New("unknown source")
		return
	}
}
