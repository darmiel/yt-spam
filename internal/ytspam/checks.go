package ytspam

import "google.golang.org/api/youtube/v3"

type Rating int

type Check interface {
}

type CommentCheck interface {
	CheckComment(comment *youtube.Comment) (Rating, error)
}

type ChannelCheck interface {
	CheckChannel(channel *youtube.Channel) (Rating, error)
}
