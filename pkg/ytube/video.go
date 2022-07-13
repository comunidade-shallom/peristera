package ytube

import (
	"html"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	VideoID     string  `json:"id"`
	Title       string  `json:"title"`
	Thumbnail   string  `json:"thumbnail"`
	PublishedAt string  `json:"publishedAt"`
	Description string  `json:"description"`
	Channel     Channel `json:"channel"`
}

type Channel struct {
	Name string `fig:"name" yaml:"name"`
	ID   string `fig:"id" yaml:"id"`
}

func FromSearchResult(raw *youtube.SearchResult) (Video, error) {
	snippet := raw.Snippet

	return Video{
		VideoID:     raw.Id.VideoId,
		Thumbnail:   snippet.Thumbnails.High.Url,
		PublishedAt: snippet.PublishedAt,
		Description: snippet.Description,
		Title:       snippet.Title,
		Channel: Channel{
			Name: snippet.ChannelTitle,
			ID:   snippet.ChannelId,
		},
	}, nil
}

func (v Video) ID() string {
	return v.VideoID
}

func (v Video) UnescapeTitle() string {
	return html.UnescapeString(v.Title)
}

func (v Video) URL() string {
	return "https://youtu.be/" + v.ID()
}
