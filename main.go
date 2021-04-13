package main

import (
	"fmt"
	"github.com/darmiel/yt-spam/pkg/ytapi"
	"github.com/darmiel/yt-spam/pkg/ytresp"
	"log"
)

func main() {
	data, err := ytapi.ScrapeChannelInfo("3blue1brown")
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}
	/*
		meta := data.Metadata.Renderer
		log.Println("Scraped channel:", meta.Title, "(", meta.ChannelURL, ")")
		log.Println("Description:", meta.Description)
		log.Println("Links:")
		for _, l := range data.GetLinks() {
			log.Println("L *", l.Title.Simpletext, ":", l.Extract())
		}
	*/
	wrap, err := ytresp.WrapYTChannelInfo(data)
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}
	fmt.Printf("%+v\n", wrap)
}
