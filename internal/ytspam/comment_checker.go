package ytspam

import (
	"encoding/base64"
	"github.com/darmiel/yt-spam/internal/checks"
	"google.golang.org/api/youtube/v3"
	"strconv"
	"time"
)

type CommentChecker struct {
	id         string
	comments   map[string]*youtube.Comment
	violations map[string][]*checks.Violation
}

func (c *CommentChecker) Violations() map[string][]*checks.Violation {
	return c.violations
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

func NewCommentChecker(comments map[string]*youtube.Comment) *CommentChecker {
	evID := newID()
	return &CommentChecker{
		id:         evID,
		comments:   comments,
		violations: make(map[string][]*checks.Violation),
	}
}

func (c *CommentChecker) addViolation(comment *youtube.Comment, check checks.Check, rating checks.Rating) {
	if c.violations == nil {
		c.violations = make(map[string][]*checks.Violation)
	}
	authorID := comment.Snippet.AuthorChannelId.Value
	violations, ok := c.violations[authorID]
	if !ok {
		violations = make([]*checks.Violation, 0)
	}
	violations = append(violations, &checks.Violation{
		Rating: rating,
		Check:  check,
	})
	c.violations[authorID] = violations
}

func (c *CommentChecker) Check(checks ...checks.CommentCheck) error {
	// clean checks
	for _, check := range checks {
		if err := check.Clean(); err != nil {
			return err
		}
	}

	for _, check := range checks {
		if err := check.CheckComments(c.comments); err != nil {
			return err
		}
	}

	// get results
	for _, check := range checks {
		for comment, rating := range check.Finalize() {
			if !rating.IsViolation() {
				continue
			}
			c.addViolation(comment, check, rating)
		}
	}
	return nil
}
