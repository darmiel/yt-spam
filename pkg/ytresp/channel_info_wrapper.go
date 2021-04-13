package ytresp

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

var numOnlyRegex *regexp.Regexp

func init() {
	var err error
	if numOnlyRegex, err = regexp.Compile(`[^0-9]`); err != nil {
		log.Fatalln(err)
	}
}

type YTChannelStats struct {
	JoinDate   time.Time // x
	TotalViews int64     // x
}

type YTChannelLink struct {
	Title string // x
	URL   string // x
}

type YTChannelInfo struct {
	URL         string            // x
	Name        string            // x
	Description string            // x
	Location    string            // x
	AvatarURLs  map[string]string // x
	BannerURLs  map[string]string // x
	Stats       *YTChannelStats   // x
	Links       []*YTChannelLink  // x
}

func WrapYTChannelInfo(data *ChannelInfoYTInitialData) (*YTChannelInfo, error) {

	res := &YTChannelInfo{
		URL:         data.Metadata.Renderer.ChannelURL,
		Name:        data.Metadata.Renderer.Title,
		Description: data.Metadata.Renderer.Description,
		BannerURLs:  make(map[string]string),
		AvatarURLs:  make(map[string]string),
	}
	stats := &YTChannelStats{}

	// get country
	for _, tab := range data.Contents.Renderer.Tabs {
		for _, c := range tab.TabRenderer.Content.Renderer.Contents {
			for _, content := range c.Renderer.Contents {
				r := content.Renderer
				// Country
				if r.Country.SimpleText != "" {
					res.Location = r.Country.SimpleText
				}
				// View
				if r.ViewcountText.SimpleText != "" {
					numOnly := numOnlyRegex.ReplaceAllString(r.ViewcountText.SimpleText, "")
					log.Println("num only:", numOnly)
					var err error
					if stats.TotalViews, err = strconv.ParseInt(numOnly, 10, 64); err != nil {
						return nil, err
					}
				}
				// Join
				for _, run := range r.JoinedDateText.Runs {
					if run.Text == "" {
						continue
					}
					if t, err := time.Parse("02.01.2006", run.Text); err == nil {
						stats.JoinDate = t
					}
				}
			}
		}
	}

	// links
	for _, link := range data.GetLinks() {
		res.Links = append(res.Links, &YTChannelLink{
			Title: link.Title.Simpletext,
			URL:   link.Extract(),
		})
	}

	// avatar / banner
	for _, banner := range data.Header.Renderer.Banner.Thumbnails {
		res.BannerURLs[fmt.Sprintf("%vx%v", banner.Width, banner.Height)] = banner.URL
	}
	for _, avatar := range data.Header.Renderer.Avatar.Thumbnails {
		res.AvatarURLs[fmt.Sprintf("%vx%v", avatar.Width, avatar.Height)] = avatar.URL
	}

	// update stats
	res.Stats = stats
	return res, nil
}
