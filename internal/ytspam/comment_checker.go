package ytspam

import (
	"encoding/base64"
	"google.golang.org/api/youtube/v3"
	"strconv"
	"time"
)

type Violation struct {
	rating Rating
	check  Check
}

type CommentChecker struct {
	id         string
	comments   map[string]*youtube.Comment
	violations map[string][]*Violation
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
		violations: make(map[string][]*Violation),
	}
}

func (c *CommentChecker) addViolation(comment *youtube.Comment, check Check, rating Rating) {
	if c.violations == nil {
		c.violations = make(map[string][]*Violation)
	}
	authorID := comment.Snippet.AuthorChannelId.Value
	violations, ok := c.violations[authorID]
	if !ok {
		violations = make([]*Violation, 0)
	}
	violations = append(violations, &Violation{
		rating: rating,
		check:  check,
	})
	//log.Println("Added Violation for Check", check.Name(),
	//	"to [", comment.Snippet.AuthorDisplayName, "] with rating", rating)
}

func (c *CommentChecker) Check(checks ...CommentCheck) error {
	// clean checks
	for _, check := range checks {
		check.Clean()
	}

	for _, comment := range c.comments {
		for _, check := range checks {
			if err := check.CheckComment(comment); err != nil {
				return err
			}
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
