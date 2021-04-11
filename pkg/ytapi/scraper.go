package ytapi

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"strings"
)

func ScrapeChannelInfo(id ChannelID) (res *channelInfoYTInitialData, err error) {
	var resp *req.Resp
	if resp, err = req.Get(id.GetChannelURL()); err != nil {
		return
	}
	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(strings.NewReader(resp.String())); err != nil {
		return
	}
	found := 0
	doc.Find("script").Each(func(i int, sel *goquery.Selection) {
		text := sel.Text()
		// we are only interested in the "ytInitialData"
		if !strings.Contains(text, "var ytInitialData =") {
			return
		}
		found++
		// parse text
		text = strings.TrimSpace(text[strings.Index(text, "=")+1:])
		text = text[:len(text)-1]
		// try to parse data
		res = new(channelInfoYTInitialData)
		err = json.Unmarshal([]byte(text), res)
	})
	if found == 0 {
		err = errors.New("not found")
	}
	return
}
