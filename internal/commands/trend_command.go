package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
)

func (cmd *Command) TrendCommand(ctx *cli.Context) error {
	max := ctx.Int64("max-videos")
	force := ctx.Bool("cache-yes")
	skipNum := ctx.Int("skip-num")
	skipIds := ctx.StringSlice("skip-ids")

	skipIdsMap := make(map[string]bool)
	for _, i := range skipIds {
		skipIdsMap[i] = true
	}

	call := cmd.Service.Videos.List([]string{"id", "snippet"}).
		Chart("mostPopular").
		RegionCode("de").
		VideoCategoryId("0").
		MaxResults(max)

	log.Println("Requesting top", max, "videos")
	resp, err := call.Do()
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("---")
	fmt.Println()
	for _, v := range resp.Items {
		log.Println("*", v.Id, ":", "'"+v.Snippet.Title+"'", "by", v.Snippet.ChannelTitle)
	}
	fmt.Println()
	fmt.Println("---")
	fmt.Println()

	i := 0
	for _, v := range resp.Items {
		i++
		log.Println(":: Video:", v.Id, "=", v.Snippet.Title, "by", v.Snippet.ChannelTitle)
		if _, ok := skipIdsMap[v.Id]; ok || (skipNum != -1 && i <= skipNum) {
			log.Println("   :: Skipped.")
			continue
		}
		if err := cmd.c(v.Id, force); err != nil {
			return err
		}
	}
	return nil
}
