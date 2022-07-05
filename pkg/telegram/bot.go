package telegram

import (
	"context"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const poolingTiming = 10 * time.Second

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
		cfg:     cfg,
		youtube: youtube,
	}

	bot.Use(useLogger(logger))

	bot.Handle("/start", handlers.Start)
	bot.Handle("/pix", handlers.Pix)
	bot.Handle("/oferta", handlers.Pix)
	bot.Handle("/sobre", handlers.Start)
	bot.Handle("/videos", handlers.Videos)

	adm := bot.Group()
	adm.Use(onlyAdmins(cfg.Telegram))

	adm.Handle("/me", handlers.Me)
	adm.Handle("/system", handlers.System)

	err = bot.SetCommands([]telebot.Command{
		{
			Text:        "sobre",
			Description: "Informações sobre a Shallom em Meriti",
		},
		{
			Text:        "oferta",
			Description: "Informações para ofertar online",
		},
		{
			Text:        "videos",
			Description: "Últimos vídeos do nosso YouTube",
		},
	})

	return bot, err
}
