package checks

import (
	"errors"
	"github.com/darmiel/yt-spam/internal/ytspam"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

var profilePicureDownladPath string

func init() {
	profilePicureDownladPath = path.Join("data", "profile-pictures")

	// check if path exists
	if _, err := os.Stat(profilePicureDownladPath); os.IsNotExist(err) {
		// make dir
		if err := os.MkdirAll(profilePicureDownladPath, 666); err != nil {
			log.Fatalln(err)
		}
	}
}

var ProfilePictureAlreadyDownloaded = errors.New("profile picture already downloaded")

type CommentDownloadProfilePictureCheck struct{}

func (c CommentDownloadProfilePictureCheck) CheckComment(comment *youtube.Comment) (ytspam.Rating, error) {
	snippet := comment.Snippet

	// check if already exists
	p := path.Join(profilePicureDownladPath, snippet.AuthorChannelId.Value+".png")
	if inf, err := os.Stat(p); inf != nil || os.IsExist(err) {
		return -1, ProfilePictureAlreadyDownloaded
	}

	// download image
	url := snippet.AuthorProfileImageUrl

	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	// read
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	// write
	if err := ioutil.WriteFile(p, b, 666); err != nil {
		return -1, err
	}

	return 0, nil
}
