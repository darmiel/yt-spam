package checks

import (
	"github.com/muesli/termenv"
	"google.golang.org/api/youtube/v3"
)

type Rating int

func (r *Rating) IsViolation() bool {
	return *r > 0
}

///

type CommentRatingNotify struct {
	Reason  string
	Comment *youtube.Comment
	Rating  Rating
	Check   NamedCheck
}

type ChannelRatingNotify struct {
	ChannelName string
	ChannelID   string
	Reason      string
	Channel     *youtube.Channel
	Rating      Rating
	Check       NamedCheck
}

///

type NamedCheck interface {
	Name() string
	Prefix() termenv.Style
	SendViolation(...interface{})
}

type SingleCommentCheck interface {
	Name() string
	Prefix() termenv.Style
	SendViolation(...interface{})

	CheckComment(comment *youtube.Comment) error
}

type MultiCommentCheck interface {
	Name() string
	Prefix() termenv.Style
	SendViolation(...interface{})

	CheckComments(comments []*youtube.Comment) error
}

type SingleChannelCheck interface {
	Name() string
	Prefix() termenv.Style
	SendViolation(...interface{})

	CheckChannel(channel *youtube.Channel) error
}

type MultiChannelCheck interface {
	Name() string
	Prefix() termenv.Style
	SendViolation(...interface{})

	CheckChannels(channels []*youtube.Channel) error
}

type CommentChannelCheck interface {
	Name() string
	Prefix() termenv.Style
	SendViolation(...interface{})

	CheckChannelByComment(comment *youtube.Comment) error
}
