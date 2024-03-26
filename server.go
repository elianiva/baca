package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

func NewServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewRenderer()

	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		source := c.QueryParam("source")
		if source == "" {
			source = "hackernews"
		}

		entries, err := fetchNews(source)
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

		entries, err := fetchNews(source)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.Render(http.StatusOK, "source-list.gohtml", map[string]any{
			"Entries": entries,
		})
	})

	e.GET("/menu/:source", func(c echo.Context) error {
		source := c.Param("source")
		c.Response().Header().Set("HX-Trigger", "list-changed")
		return c.Render(http.StatusOK, "menu-list.gohtml", map[string]any{
			"Source": source,
		})
	})

	return e
}

func fetchNews(source string) (entries []NewsEntry, err error) {
	switch source {
	case "lobsters":
		entries, err = fetchLobsters()
		return
	case "hackernews":
		entries, err = fetchHackernews()
		return
	default:
		err = errors.New("unknown source")
		return
	}
}
