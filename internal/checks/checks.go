package checks

import (
	"google.golang.org/api/youtube/v3"
)

type Rating int

func (r *Rating) IsViolation() bool {
	return *r > 0
}

type Check interface {
	Name() string
	Finalize() map[*youtube.Comment]Rating
	Clean() error
}

type CommentCheck interface {
	Name() string
	Finalize() map[*youtube.Comment]Rating
	Clean() error
	CheckComments(all map[string]*youtube.Comment) error
}

type ChannelCheck interface {
	Name() string
	Finalize() map[*youtube.Comment]Rating
	Clean() error
	CheckChannels(all []*youtube.Channel) error
}
