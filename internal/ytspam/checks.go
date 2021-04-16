package ytspam

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
	Clean()
}

type CommentCheck interface {
	Name() string
	Finalize() map[*youtube.Comment]Rating
	Clean()
	CheckComment(comment *youtube.Comment) error
}

type ChannelCheck interface {
	Name() string
	Finalize() map[*youtube.Comment]Rating
	Clean()
	CheckChannel(channel *youtube.Channel) error
}
