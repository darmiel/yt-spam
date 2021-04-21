package ytspam

import (
	"encoding/json"
	"github.com/cheggaaa/pb/v3"
	"google.golang.org/api/youtube/v3"
	"log"
)

type CommentReaderConfig struct {
	DisplayBar bool
	Silent     bool
}

type CommentReader struct {
	service  *youtube.Service
	video    *youtube.Video
	comments map[string]*youtube.Comment
	stats    *CommentReaderStats
	bar      *pb.ProgressBar
	*CommentReaderConfig
}

type CommentReaderStats struct {
	ReadComments int64
	ReadReplies  int64
}

func NewCommentReader(service *youtube.Service, video *youtube.Video, config ...*CommentReaderConfig) *CommentReader {
	var cfg *CommentReaderConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = &CommentReaderConfig{}
	}
	return &CommentReader{
		service:             service,
		video:               video,
		CommentReaderConfig: cfg,
	}
}

func (r *CommentReader) FromJSON(data []byte) (err error) {
	var val []*youtube.Comment
	if err = json.Unmarshal(data, &val); err != nil {
		return
	}
	for _, c := range val {
		r.addComment(c)
	}
	return
}

func (r *CommentReader) ToJSON() (res []byte, err error) {
	var val []*youtube.Comment
	for _, c := range r.comments {
		val = append(val, c)
	}
	res, err = json.Marshal(val)
	return
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

func (r *CommentReader) readCommentReplies(comment *youtube.Comment, token ...string) error {
	if r.stats == nil {
		r.stats = &CommentReaderStats{}
	}
	call := r.service.Comments.List([]string{"id", "snippet"}).ParentId(comment.Id)
	if len(token) > 0 {
		call = call.PageToken(token[0])
	}
	resp, err := call.Do()
	if err != nil {
		return err
	}
	r.stats.ReadReplies += int64(len(resp.Items))
	for _, re := range resp.Items {
		r.addComment(re)
	}
	if resp.NextPageToken != "" {
		return r.readCommentReplies(comment, resp.NextPageToken)
	}
	return nil
}

func (r *CommentReader) readVideoComments(npt ...string) error {
	if r.stats == nil {
		r.stats = &CommentReaderStats{}
	}
	call := r.service.CommentThreads.List([]string{"id", "replies", "snippet"}).
		// Filters
		VideoId(r.video.Id).
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

	r.stats.ReadComments += int64(len(resp.Items))
	for _, t := range resp.Items {
		comment := t.Snippet.TopLevelComment
		r.addComment(comment)

		// check replies
		if t.Replies != nil {
			repl := t.Replies.Comments
			if int64(len(t.Replies.Comments)) < t.Snippet.TotalReplyCount {
				if err = r.readCommentReplies(comment); err != nil {
					return err
				}
			} else {
				r.stats.ReadReplies += int64(len(repl))
				for _, reply := range repl {
					r.addComment(reply)
				}
			}
		}
	}

	if r.DisplayBar {
		if r.bar == nil {
			r.bar = pb.Full.Start64(int64(r.video.Statistics.CommentCount))
		}
		r.bar.SetCurrent(int64(len(r.comments)))
	} else if !r.Silent {
		log.Println("Read", r.stats.ReadComments, "comments +",
			r.stats.ReadReplies, "replies =", r.stats.ReadReplies+r.stats.ReadComments,
			"(", len(r.comments), ")")
	}

	if resp.NextPageToken != "" {
		return r.readVideoComments(resp.NextPageToken)
	}

	return nil
}

func (r *CommentReader) Read() error {
	// clear old comments
	r.comments = make(map[string]*youtube.Comment)
	defer func() {
		if r.bar != nil {
			r.bar.Finish()
			log.Println("finished bar.")
		}
	}()
	return r.readVideoComments()
}
