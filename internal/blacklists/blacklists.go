package blacklists

import (
	"github.com/darmiel/yt-spam/internal/compare"
	"log"
	"path"
)

// blacklists

// channel
var (
	channelBlacklistFile = "channel-blacklist.txt"
	ChannelBlacklist     = mustRead(channelBlacklistFile)
)

// comment
var (
	commentBlacklistFile = "comment-blacklist.txt"
	CommentBlacklist     = mustRead(commentBlacklistFile)
)

// name
var (
	nameBlacklistFile = "name-blacklist.txt"
	NameBlacklist     = mustRead(nameBlacklistFile)
)

// playlist
var (
	playlistBlacklistFile = "playlist-blacklist.txt"
	PlaylistBlacklist     = mustRead(playlistBlacklistFile)
)

// other

var (
	InputDataPath = path.Join("data", "input")
)

func fatal(err error, file string) {
	var msg string
	if err != nil {
		msg = ": " + err.Error()
	} else {
		msg = ""
	}
	log.Fatal("FATAL :: [Blacklist] Failed to open file '", file, "'", msg, "\n")
}

func mustRead(file string) []compare.StringCompare {
	data, err := compare.FromFile(path.Join(InputDataPath, file))
	if err != nil {
		fatal(err, channelBlacklistFile)
		return nil
	}
	return data
}
