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
		Title: video.Title(),
	}, nil
}

func (v Video) Hash() string {
	hash := sha256.Sum256([]byte(v.URL))

	return base64.RawStdEncoding.EncodeToString(hash[:])
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
