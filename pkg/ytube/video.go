package ytube

import (
	"encoding/json"
	"time"

	"google.golang.org/api/youtube/v3"
)

type VideoID struct {
	VideoID string `json:"videoId"`
}

type VideoThumbnails struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

type VideoSnippet struct {
	Title        string                     `json:"title"`
	ChannelTitle string                     `json:"channelTitle"`
	PublishedAt  time.Time                  `json:"publishedAt"`
	Description  time.Time                  `json:"description"`
	Thumbnails   map[string]VideoThumbnails `json:"thumbnails"`
}

type Video struct {
	VideoID VideoID      `json:"id"`
	Snippet VideoSnippet `json:"snippet"`
}

type Channel struct {
	Name string `fig:"name" yaml:"name"`
	ID   string `fig:"id" yaml:"id"`
}

func FromSearchResult(raw *youtube.SearchResult) (Video, error) {
	var res Video

	bt, err := raw.MarshalJSON()
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bt, &res)

	return res, err
}

func (v Video) URL() string {
	return "https://youtu.be/" + v.VideoID.VideoID
}

func (v Video) ToBotContent() (interface{}, error) {
	return v.URL(), nil
}
