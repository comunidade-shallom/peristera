package cron

import (
	"context"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type MessageBuilder func() ([]interface{}, error)

type Sendable interface {
	ToBotContent() (interface{}, error)
}

type Jobs struct {
	cfg       config.AppConfig
	bot       *telebot.Bot
	youtube   ytube.Service
	scheduler *gocron.Scheduler
}

var ErrNoChatToBroadcast = errors.Business("No chats to broadcast", "CRON:001")

func (j Jobs) Start(ctx context.Context) error {
	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		err := recoverData.(error)

		zerolog.Ctx(ctx).
			Error().
			Str("context", "jobs:panic-handler").
			Str("job-name", jobName).
			Err(err).
			Stack().
			Msg("Panic to run job")
	})

	j.scheduler.StartAsync()

	<-ctx.Done()

	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "jobs:manager").
		Logger()

	logger.Warn().Msg("Stopping cron jobs...")

	j.scheduler.Stop()

	logger.Warn().Msg("Cron jobs are stopped")

	return ctx.Err()
}

func (j Jobs) jobLogger(ctx context.Context, job string) zerolog.Logger {
	return zerolog.Ctx(ctx).
		With().
		Str("context", "jobs").
		Str("job", job).
		Logger()
}

func (j *Jobs) register(ctx context.Context) error {
	timezone, err := time.LoadLocation(j.cfg.Timezone)
	if err != nil {
		return err
	}

	j.scheduler = gocron.NewScheduler(timezone)

	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "jobs:register").
		Logger()

	if len(j.cfg.Cron.LastUpdates) == 0 {
		logger.Warn().Msg("No cron entries from last updates")

		return nil
	}

	logger.Info().Msg("Registring last updates")

	lastUpdates := func() error {
		return j.LastVideos(ctx)
	}

	j.scheduler.SingletonModeAll()

	for _, entry := range j.cfg.Cron.LastUpdates {
		_, err := j.scheduler.Cron(entry).Do(lastUpdates)
		if err != nil {
			return err
		}

		logger.Info().Msgf("Registred: %s", entry)
	}

	return nil
}

func (j Jobs) sendMessages(chatId int64, entries []interface{}, opts ...interface{}) error {
	var err error

	for _, entry := range entries {
		content, _ := toBotContent(entry)

		_, err = j.bot.Send(&telebot.Chat{
			ID: chatId,
		}, content, opts...)

		if err != nil {
			return err
		}
	}

	return nil
}

func toBotContent(source interface{}) (interface{}, error) {
	switch val := source.(type) {
	case Sendable:
		return val.ToBotContent()
	default:
		return source, nil
	}
}

func (j Jobs) broadcast(fn MessageBuilder, opts ...interface{}) error {
	if len(j.cfg.Telegram.Broadcast) == 0 {
		return ErrNoChatToBroadcast
	}

	entries, err := fn()
	if err != nil {
		return err
	}

	for _, chatId := range j.cfg.Telegram.Broadcast {
		err = j.sendMessages(chatId, entries, opts...)

		if err != nil {
			return err
		}
	}

	return nil
}
