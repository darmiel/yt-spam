package ytspam

import (
	"encoding/base64"
	"fmt"
	"github.com/darmiel/yt-spam/internal/checks"
	"google.golang.org/api/youtube/v3"
	"io"
	"strconv"
	"sync"
	"time"
)

type CommentChecker struct {
	id       string
	comments []*youtube.Comment

	ResErrors   chan interface{}
	ResComments chan *checks.CommentRatingNotify
	ResChannels chan *checks.ChannelRatingNotify
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
		id:          evID,
		comments:    comments,
		ResErrors:   make(chan interface{}),
		ResComments: make(chan *checks.CommentRatingNotify),
		ResChannels: make(chan *checks.ChannelRatingNotify),
	}
}

func (c *CommentChecker) executeSingle(wg *sync.WaitGroup, fc func(comment *youtube.Comment)) {
	for _, comment := range c.comments {
		fc(comment)
	}
	if err := recover(); err != nil {
		c.ResErrors <- err
	}
	wg.Done()
}

func (c *CommentChecker) executeMulti(wg *sync.WaitGroup, fc func(comments []*youtube.Comment)) {
	fc(c.comments)
	if err := recover(); err != nil {
		c.ResErrors <- err
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

	close(c.ResChannels)
	close(c.ResComments)
	close(c.ResErrors)
}

func (c *CommentChecker) ReadComments(w io.Writer, wg *sync.WaitGroup) {
	for notify := range c.ResComments {
		var msg string
		if notify.Reason != "" {
			msg = notify.Reason
		} else {
			var id string
			if notify.Comment != nil {
				id = notify.Comment.Id
			} else {
				id = "unknown"
			}
			msg = fmt.Sprintf("[Comment/%s] r = %v; ID: %s",
				notify.Check.Name(), notify.Rating, id)
		}
		_, _ = fmt.Fprintln(w, msg)
	}
	wg.Done()
}

func (c *CommentChecker) ReadChannels(w io.Writer, wg *sync.WaitGroup) {
	for notify := range c.ResChannels {
		var msg string
		if notify.Reason != "" {
			msg = notify.Reason
		} else {
			msg = fmt.Sprintf("[Channel/%s] r = %v; %s (%s)",
				notify.Check, notify.Rating, notify.ChannelName, notify.ChannelID)
		}
		_, _ = fmt.Fprintln(w, msg)
	}
	wg.Done()
}

func (c *CommentChecker) ReadErrors(w io.Writer, wg *sync.WaitGroup) {
	for notify := range c.ResErrors {
		_, _ = fmt.Fprintf(w, "ERROR: %v\n", notify)
	}
	wg.Done()
}
