package copycat

import (
	"google.golang.org/api/youtube/v3"
	"log"
)

const (
	CopyCatMinCommentLength = 24
)

func (c *CommentCopyCatCheck) CheckComments(comments map[string]*youtube.Comment) error {
	checked := make(map[string]bool)

	for _, cursor := range comments {
		// check length
		if len(trimBody(cursor)) < CopyCatMinCommentLength {
			continue
		}

		if _, checked := checked[cursor.Id]; checked {
			continue
		}

		earliest := cursor
		matches := getMatchingComments(cursor, comments)
		for _, matching := range matches {
			checked[matching.Id] = true
			if commentBefore(matching, earliest) {
				earliest = matching
			}
		}

		for _, matching := range matches {
			if matching != earliest {
				// check if author copied himself
				if commentFromHimHimself(matching, earliest) {
					continue
				}
				c.printCCMessage(earliest, matching)

				// add violation
				c.violations[matching] = 25
			}
		}
	}

	log.Println("Checked a total of", len(checked), "comments.")
	return nil
}
