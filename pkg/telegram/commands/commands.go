package commands

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Commands struct {
	cfg     config.AppConfig
	youtube ytube.Service
}

func New(cfg config.AppConfig, youtube ytube.Service) Commands {
	return Commands{
		cfg:     cfg,
		youtube: youtube,
	}
}

func (h Commands) Setup(ctx context.Context, bot *telebot.Bot) error {
	logger := zerolog.Ctx(ctx).With().Str("context", "commands").Logger()

	if h.cfg.Debug {
		bot.Use(middleware.Logger())
	}

	bot.Use(useLogger(logger))

	bot.Handle("/start", h.Start)
	bot.Handle("/pix", h.Pix)
	bot.Handle("/oferta", h.Pix)
	bot.Handle("/sobre", h.Start)
	bot.Handle("/videos", h.Videos)

	adm := bot.Group()
	adm.Use(onlyAdmins(h.cfg.Telegram))

	adm.Handle("/me", h.Me)
	adm.Handle("/system", h.System)

	return bot.SetCommands([]telebot.Command{
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
}

func (h Commands) logger(tx telebot.Context) zerolog.Logger {
	return tx.Get(loggerKey).(zerolog.Logger) //nolint:forcetypeassert
}