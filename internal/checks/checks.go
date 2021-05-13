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

// ok
type SingleCommentCheck interface {
	NamedCheck
	CheckComment(comment *youtube.Comment)
}

type MultiCommentCheck interface {
	NamedCheck
	CheckComments(comments []*youtube.Comment)
}

type SingleChannelCheck interface {
	NamedCheck
	CheckChannel(channel *youtube.Channel)
}

type MultiChannelCheck interface {
	NamedCheck
	CheckChannels(channels []*youtube.Channel)
}

type CommentChannelCheck interface {
	NamedCheck
	CheckChannelByComment(comment *youtube.Comment)
}
