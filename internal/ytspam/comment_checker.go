package ytspam

import (
	"encoding/base64"
	"github.com/darmiel/yt-spam/internal/checks"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/errgo.v2/fmt/errors"
	"strconv"
	"sync"
	"time"
)

type CommentChecker struct {
	id       string
	comments []*youtube.Comment

	resErrors   chan error
	resComments chan *checks.CommentRatingNotify
	resChannels chan *checks.ChannelRatingNotify
}

// CurrentTime << 4 | Worker
var evNum int64

func newID() string {
	var val = (time.Now().Unix() / 1000) << 4
	evNum++
	val |= evNum
	b := []byte(strconv.FormatInt(val, 10))
	return base64.StdEncoding.EncodeToString(b)
}

///

func NewCommentChecker(comments []*youtube.Comment) *CommentChecker {
	evID := newID()
	return &CommentChecker{
		id:       evID,
		comments: comments,
	}
}

func (c *CommentChecker) executeSingle(wg *sync.WaitGroup, fc func(comment *youtube.Comment)) {
	for _, comment := range c.comments {
		fc(comment)
	}
	if err := recover(); err != nil {
		c.resErrors <- errors.Newf("%v", err)
	}
	wg.Done()
}

func (c *CommentChecker) executeMulti(wg *sync.WaitGroup, fc func(comments []*youtube.Comment)) {
	fc(c.comments)
	if err := recover(); err != nil {
		c.resErrors <- errors.Newf("%v", err)
	}
	wg.Done()
}

func (c *CommentChecker) Check(ch ...checks.NamedCheck) {
	var wg sync.WaitGroup
	for _, check := range ch {
		// 1. Channel Checks
		if check, ok := check.(checks.CommentChannelCheck); ok {
			wg.Add(1)
			go c.executeSingle(&wg, check.CheckChannelByComment)
		}

		// 2. Single Comments
		if check, ok := check.(checks.SingleCommentCheck); ok {
			wg.Add(1)
			go c.executeSingle(&wg, check.CheckComment)
		}

		// 3. Multiple Comments
		if check, ok := check.(checks.MultiCommentCheck); ok {
			wg.Add(1)
			go c.executeMulti(&wg, check.CheckComments)
		}
	}
	wg.Wait()
}
