package copycat

import (
	"fmt"
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

		if len(matches) > 4 {
			fmt.Println()
			log.Println("### ", cursor.Snippet.AuthorDisplayName, ":", trimBody(cursor), ":: copied", len(matches)-1, "times")
		}
		for _, matching := range matches {
			if matching != earliest {
				// check if author copied himself
				if commentFromHimHimself(matching, earliest) {
					continue
				}
				printCCMessage(earliest, matching)

				// add violation
				c.violations[matching] = 25
			}
		}
		if len(matches) > 4 {
			fmt.Println()
		}
	}

	log.Println("Checked a total of", len(checked), "comments.")
	return nil
}
