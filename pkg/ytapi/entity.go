package ytapi

import "fmt"

type ChannelID string

func (c *ChannelID) GetChannelURL() string {
	return fmt.Sprintf("https://www.youtube.com/c/%s/about", string(*c))
}
