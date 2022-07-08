package ytube

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/rs/zerolog"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	youtube *youtube.Service
}

func NewService(ctx context.Context, cfg config.YouTube) (Service, error) {
	yService, err := youtube.NewService(ctx, option.WithAPIKey(cfg.Token))
	if err != nil {
		return Service{}, err
	}

	return Service{
		youtube: yService,
	}, nil
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

	logger := zerolog.Ctx(ctx).
		With().
		Str("fn", "youtube:LastVideos").
		Logger()

	for _, item := range res.Items {
		vid, err := FromSearchResult(item)
		if err != nil {
			logger.Warn().Err(err).Msg("Parse error: %s")

			continue
		}

		vids = append(vids, vid)
	}

	return vids, err
}
