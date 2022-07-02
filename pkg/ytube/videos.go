package ytube

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/pterm/pterm"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func NewService(ctx context.Context, cfg config.AppConfig) (Service, error) {
	yService, err := youtube.NewService(ctx, option.WithAPIKey(cfg.YoutubeToken))
	if err != nil {
		return Service{}, err
	}

	return Service{
		youtube: yService,
	}, nil
}

type Service struct {
	youtube *youtube.Service
}

func (s Service) LastVideos(ctx context.Context, channelID string, maxResults int) ([]Video, error) {
	res, err := s.youtube.Search.
		List([]string{"snippet"}).
		MaxResults(int64(maxResults)).
		Order("date").
		ChannelId(channelID).
		Context(ctx).
		Do()

	vids := []Video{}

	if err != nil {
		return vids, err
	}

	for _, item := range res.Items {
		vid, err := FromSearchResult(item)

		if err == nil {
			pterm.Warning.Printfln("Parse error: %s", err.Error())

			continue
		}

		vids = append(vids, vid)
	}

	return vids, err
}
