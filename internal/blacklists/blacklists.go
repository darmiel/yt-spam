package blacklists

import (
	"github.com/darmiel/yt-spam/internal/compare"
	"log"
	"path"
)

// blacklists
type StringBlacklist []compare.StringCompare

func (b StringBlacklist) AnyMatch(s string) compare.StringCompare {
	for _, c := range b {
		if c.Compare(s) {
			return c
		}
	}
	return nil
}

func (b StringBlacklist) AnyAnyMatch(s ...string) compare.StringCompare {
	for i, str := range s {
		if i > 0 && str == s[i-1] {
			continue
		}
		if c := b.AnyMatch(str); c != nil {
			return c
		}
	}
	return nil
}

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

func mustRead(file string) StringBlacklist {
	data, err := compare.FromFile(path.Join(InputDataPath, file))
	if err != nil {
		fatal(err, channelBlacklistFile)
		return nil
	}
	return data
}
