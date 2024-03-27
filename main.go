package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Name:  "baca",
		Usage: "A minimalist news aggregator from various sources",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Start the application server",
				Action: func(c *cli.Context) error {
					cache := newCache()
					server := newServer(&ServerConfig{
						Cache: cache,
					})

					exitSig := make(chan os.Signal, 1)
					signal.Notify(exitSig, os.Interrupt, os.Kill)

					go func() {
						<-exitSig
						ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
						defer cancel()

						if err := server.Shutdown(ctx); err != nil {
							log.Error().Err(err).Msg("failed to shutdown server")
						}
					}()

					err := server.Start(":8080")
					if err != nil && !errors.Is(err, http.ErrServerClosed) {
						log.Fatal().Err(err).Msg("failed to start server")
					}

					return nil
				},
			},
		},
	}
}

func main() {
	if err := App().Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run app")
	}
}
