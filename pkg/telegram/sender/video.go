package sender

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/comunidade-shallom/peristera/pkg/ytube"
)

type Video struct {
	Chats
	URL   string
	Title string
}

func FromVideo(video ytube.Video, chats Chats) (Video, error) {
	return Video{
		Chats: chats,
		URL:   video.URL(),
		Title: video.UnescapeTitle(),
	}, nil
}

func (v Video) Hash() []byte {
	raw := sha256.Sum256([]byte(v.URL))
	hash := raw[:]

	buf := make([]byte, base64.RawStdEncoding.EncodedLen(len(hash)))

	base64.RawStdEncoding.Encode(buf, hash)

	return append([]byte("videos:"), buf...)
}

func (v Video) Params() []interface{} {
	return []interface{}{}
}

func (v Video) Content() interface{} {
	var buider strings.Builder

	buider.WriteString(v.Title)
	buider.WriteRune('\n')
	buider.WriteString(v.URL)

	return buider.String()
}
