package cron

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/telegram/sender"
)

const videosMax = 2

func (j *Jobs) LastVideos(ctx context.Context) error {
	logger := j.jobLogger(ctx, "last-videos")

	logger.Info().Msg("Running...")

	for _, channel := range j.cfg.Youtube.Channels {
		err := j.lastChannelVideos(ctx, channel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *Jobs) lastChannelVideos(ctx context.Context, channel config.Channel) error {
	return j.broadcast(ctx, func(chats sender.Chats) ([]sender.Sendable, error) {
		logger := j.jobLogger(ctx, "last-videos:"+channel.ID)

		logger.Info().Msg("Loading last videos...")

		vids, err := j.youtube.LastVideos(ctx, channel.ID, videosMax)

		length := len(vids)
		logger.Info().Msgf("Videos loaded: %v", length)

		if err != nil {
			logger.Error().Err(err).Msg("Fail to load last videos")

			return make([]sender.Sendable, 0), err
		}

		res := make([]sender.Sendable, length)

		for index, vid := range vids {
			msg, err := sender.FromVideo(vid, chats)
			if err != nil {
				return make([]sender.Sendable, 0), err
			}

			logger.Info().
				Str("id", vid.ID()).
				Str("title", vid.Title).
				Msgf("Video loaded")

			res[index] = msg
		}

		return res, nil
	})
}
