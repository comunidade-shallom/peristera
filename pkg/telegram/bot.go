package telegram

import (
	"context"
	"strconv"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const poolingTiming = 10 * time.Second

const loggerKey = "logger"

func NewBot(ctx context.Context, cfg config.AppConfig) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: poolingTiming},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return bot, err
	}

	youtube, err := ytube.NewService(ctx, cfg.Youtube)
	if err != nil {
		return bot, err
	}

	if cfg.Debug {
		bot.Use(middleware.Logger())
	}

	logger := zerolog.Ctx(ctx).With().Str("context", "bot").Logger()

	handlers := Handler{
		logger:  logger,
		cfg:     cfg,
		youtube: youtube,
	}

	bot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			c.Set(
				loggerKey,
				logger.With().
					Str("message_id", strconv.Itoa(c.Message().ID)).
					Str("sender_id", strconv.FormatInt(c.Sender().ID, 10)). //nolint:gomnd
					Str("chat_id", strconv.FormatInt(c.Chat().ID, 10)).     //nolint:gomnd
					Logger(),
			)

			return next(c)
		}
	})

	bot.Handle("/me", handlers.Me)
	bot.Handle("/videos", handlers.Videos)

	return bot, nil
}
