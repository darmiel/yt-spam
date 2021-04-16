package ytspam

import (
	"google.golang.org/api/youtube/v3"
	"log"
)

type CommentReader struct {
	VideoID  string
	comments map[string]*youtube.Comment
	service  *youtube.Service
}

func NewCommentReader(service *youtube.Service, videoID string) *CommentReader {
	return &CommentReader{
		VideoID: videoID,
		service: service,
	}
}

//////////////////////////////////////////////////////////////////////////////////////////

func (r *CommentReader) addComment(comment *youtube.Comment) {
	if comment == nil {
		return
	}
	if r.comments == nil {
		r.comments = make(map[string]*youtube.Comment)
	} else {
		// skip duplicates
		if _, ok := r.comments[comment.Id]; ok {
			return
		}
	}
	r.comments[comment.Id] = comment
}

func (r *CommentReader) GetComments() map[string]*youtube.Comment {
	return r.comments
}

func (r *CommentReader) readVideoComments(npt ...string) error {
	call := r.service.CommentThreads.List([]string{"id", "replies", "snippet"}).
		// Filters
		VideoId(r.VideoID).
		// Parameters
		MaxResults(100).
		Order("time")

	if len(npt) > 0 {
		call = call.PageToken(npt[0])
	}

	resp, err := call.Do()
	if err != nil {
		return err
	}

	for _, t := range resp.Items {
		comment := t.Snippet.TopLevelComment

		// check tlc
		r.addComment(comment)
		// check replies
		if t.Replies != nil {
			for _, reply := range t.Replies.Comments {
				r.addComment(reply)
			}
		}
	}

	// TODO: remove debug
	log.Println("Read", len(r.comments), "comments...")

	if resp.NextPageToken != "" {
		return r.readVideoComments(resp.NextPageToken)
	}

	return nil
}

func (r *CommentReader) Read() error {
	// clear old comments
	r.comments = make(map[string]*youtube.Comment)
	return r.readVideoComments()
}
