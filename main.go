package main

import (
	"github.com/darmiel/yt-spam/internal/commands"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	// read api key from "api-key.txt"
	b, err := ioutil.ReadFile("api-key.txt")
	if err != nil {
		log.Fatalln("Error reading api-key.txt:", err)
		return
	}

	apiKey := string(b)

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalln("ERROR:", err)
		return
	}

	cmd := &commands.Command{Service: service}

	app := &cli.App{
		Authors: []*cli.Author{
			{
				Name:  "darmiel",
				Email: "hi@d2a.io",
			},
		},
		Version:     "0.1.0",
		Description: "Anti YT-Spam",
		Commands: []*cli.Command{
			{
				Name:        "check",
				Aliases:     []string{"c"},
				UsageText:   "applies all checks on a video",
				Description: "applies all checks on a video",
				Category:    "checking",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Required: true,
						Name:     "video-id",
						Aliases:  []string{"i"},
					},
					&cli.BoolFlag{Name: "cache-yes"},
				},
				Action: cmd.CheckCommand,
			},
			{
				Name:        "trends",
				Aliases:     []string{"t"},
				UsageText:   "get videos from trends",
				Description: "get videos from trends",
				Category:    "checking",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:    "max-videos",
						Aliases: []string{"m"},
						Value:   25,
					},
					&cli.BoolFlag{Name: "cache-yes"},
					&cli.IntFlag{Name: "skip"},
				},
				Action: cmd.TrendCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
