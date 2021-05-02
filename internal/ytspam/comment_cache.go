package ytspam

import (
	"encoding/json"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"path"
	"time"
)

type CachedComments struct {
	VideoID    string             `json:"video_id"`
	VideoTitle string             `json:"video_title"`
	LastUpdate time.Time          `json:"last_update"`
	Comments   []*youtube.Comment `json:"comments"`
}

func (c *CachedComments) Save() error {
	p := path.Join("data", "cache", c.VideoID+".json")
	return c.SaveTo(p)
}

func (c *CachedComments) SaveTo(p string) error {
	dat, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p, dat, 0766)
}

func (c *CachedComments) Wrap() map[string]*youtube.Comment {
	res := make(map[string]*youtube.Comment)
	for _, c := range c.Comments {
		res[c.Id] = c
	}
	return res
}

func CacheFromPath(p string) (*CachedComments, error) {
	dat, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	c := new(CachedComments)
	if err := json.Unmarshal(dat, c); err != nil {
		return nil, err
	}
	return c, nil
}

func CacheFromVideoID(videoID string) (*CachedComments, error) {
	p := path.Join("data", "cache", videoID+".json")
	return CacheFromPath(p)
}
