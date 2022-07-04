package telegram

import (
	"context"
	"strconv"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const poolingTiming = 10 * time.Second

const loggerKey = "logger"

func NewBot(ctx context.Context, cfg config.AppConfig, youtube ytube.Service) (*telebot.Bot, error) {
	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: poolingTiming},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return bot, err
	}

	if cfg.Debug {
		bot.Use(middleware.Logger())
	}

	logger := zerolog.Ctx(ctx).With().Str("context", "bot").Logger()

	handlers := Handler{
		bot:     bot,
		logger:  logger,
		cfg:     cfg,
		youtube: youtube,
	}

	bot.Use(func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(tx telebot.Context) error {
			tx.Set(
				loggerKey,
				logger.With().
					Str("message_id", strconv.Itoa(tx.Message().ID)).
					Str("sender_id", strconv.FormatInt(tx.Sender().ID, 10)). //nolint:gomnd
					Str("chat_id", strconv.FormatInt(tx.Chat().ID, 10)).     //nolint:gomnd
					Logger(),
			)

			return next(tx)
		}
	})

	bot.Handle("/start", handlers.Start)
	bot.Handle("/me", handlers.Me)
	bot.Handle("/pix", handlers.Pix)
	bot.Handle("/oferta", handlers.Pix)
	bot.Handle("/videos", handlers.Videos)

	err = bot.SetCommands([]telebot.Command{
		{
			Text:        "oferta",
			Description: "Informações para ofertar online",
		},
		{
			Text:        "videos",
			Description: "Últimos vídeos do nosso youtube",
		},
	})

	return bot, err
}
