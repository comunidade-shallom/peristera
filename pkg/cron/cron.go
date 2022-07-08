package cron

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/comunidade-shallom/peristera/pkg/telegram/sender"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type MessageBuilder func(sender.Chats) ([]sender.Sendable, error)

type Options struct {
	Bot      *telebot.Bot
	Database database.Database
	YouTube  ytube.Service
	Config   config.AppConfig
}

type Jobs struct {
	senderCh  chan sender.Sendable
	cfg       config.AppConfig
	youtube   ytube.Service
	database  database.Database
	bot       *telebot.Bot
	scheduler *gocron.Scheduler
}

var ErrNoChatToBroadcast = errors.Business("No chats to broadcast", "CRON:001")

func New(ctx context.Context, opts Options) (*Jobs, error) {
	jb := &Jobs{
		youtube:  opts.YouTube,
		bot:      opts.Bot,
		cfg:      opts.Config,
		database: opts.Database,
	}

	return jb, jb.register(ctx)
}

func (j *Jobs) Start(ctx context.Context) error {
	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		err, _ := recoverData.(error)

		zerolog.Ctx(ctx).
			Error().
			Str("context", "jobs:panic-handler").
			Str("job-name", jobName).
			Err(err).
			Stack().
			Msg("Panic to run job")
	})

	j.senderCh = make(chan sender.Sendable, 2) //nolint:gomnd

	var wg sync.WaitGroup

	wg.Add(1) // SendableWorker

	go func() {
		sender.SendableWorker(ctx, j.senderCh, j.bot, j.database)
		wg.Done()
	}()

	j.scheduler.StartAsync()

	<-ctx.Done()

	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "jobs:manager").
		Logger()

	logger.Warn().Err(ctx.Err()).Msg("Stopping cron jobs...")

	j.scheduler.Stop()

	close(j.senderCh)

	logger.Warn().Msg("Cron jobs are stopped")

	wg.Wait()

	logger.Warn().Msg("Waiting for workers...")

	return nil
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
		_, err := j.scheduler.
			Cron(entry).
			Do(lastUpdates)
		if err != nil {
			return err
		}

		logger.Info().Msgf("Registred: %s", entry)
	}

	return nil
}

func (j *Jobs) jobLogger(ctx context.Context, job string) zerolog.Logger {
	return zerolog.Ctx(ctx).
		With().
		Str("context", "jobs").
		Str("job", job).
		Logger()
}

func (j *Jobs) broadcast(ctx context.Context, fn MessageBuilder, opts ...interface{}) error {
	if len(j.cfg.Telegram.Broadcast) == 0 {
		return ErrNoChatToBroadcast
	}

	entries, err := fn(j.cfg.Telegram.Broadcast)
	if err != nil {
		return err
	}

	entries, err = j.filter(ctx, entries)

	if err != nil {
		return err
	}

	for _, msg := range entries {
		j.senderCh <- msg
	}

	return nil
}

func (j *Jobs) filter(ctx context.Context, in []sender.Sendable) ([]sender.Sendable, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "jobs:message-filter").
		Logger()

	keys := make([][]byte, len(in))

	for index, m := range in {
		keys[index] = m.Hash()
	}

	logger.Info().Msg("Filtering messages...")

	keys, err := j.database.MissingKeys(ctx, keys)
	if err != nil {
		return []sender.Sendable{}, err
	}

	if len(keys) == 0 {
		logger.Warn().Msg("No messages to send")

		return []sender.Sendable{}, nil
	}

	out := make([]sender.Sendable, len(keys))

	for index, msg := range in {
		found := false
		hash := msg.Hash()

		for _, k := range keys {
			if bytes.Equal(hash, k) {
				out[index] = msg
				found = true

				continue
			}
		}

		if !found {
			logger.Warn().Bytes("hash", hash).Msg("Message is already sent")
		}
	}

	logger.Info().Msgf("Messages able to be send: %v", len(out))

	return out, nil
}
