package commands

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Commands struct {
	cfg     config.AppConfig
	youtube ytube.Service
	db      database.Database
}

func New(cfg config.AppConfig, youtube ytube.Service, db database.Database) Commands {
	return Commands{
		db:      db,
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

	_ = h.registerMenu(bot)

	bot.Handle("/start", h.Start)
	bot.Handle("/sobre", h.Start)
	bot.Handle("/pix", h.Pix)
	bot.Handle("/oferta", h.Pix)
	bot.Handle("/videos", h.Videos)
	bot.Handle("/address", h.Address)
	bot.Handle("/location", h.Address)
	bot.Handle("/endereco", h.Address)
	bot.Handle("/endereço", h.Address)
	bot.Handle("/agenda", h.Calendar)
	bot.Handle("/calendar", h.Calendar)

	adm := bot.Group()
	adm.Use(restrictTo(h.cfg.Telegram.Admins, "admins"))

	adm.Handle("/me", h.Me)
	adm.Handle("/system", h.System)

	root := bot.Group()
	root.Use(restrictTo(h.cfg.Telegram.Roots, "roots"))
	root.Handle("/exec", h.Exec)
	root.Handle("/backup", h.Backup)
	root.Handle("/load", h.Load)

	return bot.SetCommands([]telebot.Command{
		{
			Text:        "sobre",
			Description: "Informações sobre a Shallom em Meriti",
		},
		{
			Text:        "endereco",
			Description: "Nosso endereço",
		},
		{
			Text:        "agenda",
			Description: "Nossos horários de culto",
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

func (h Commands) registerMenu(bot *telebot.Bot) *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}

	bot.Use(useMenu(menu))

	btnAbout := menu.Text("ℹ️ Sobre")
	btnAgenda := menu.Text("🗓️ Agenda")
	btnAddress := menu.Text("📍 Endereço")
	btnPix := menu.Text("🏦 Pix")
	btnYoutube := menu.Text("📹 YouTube")

	menu.Reply(
		menu.Row(btnAbout, btnAddress),
		menu.Row(btnAgenda, btnPix),
		menu.Row(btnYoutube),
	)

	bot.Handle(&btnAgenda, h.Calendar)
	bot.Handle(&btnAddress, h.Address)
	bot.Handle(&btnAbout, h.Start)
	bot.Handle(&btnPix, h.Pix)
	bot.Handle(&btnYoutube, h.Videos)

	return menu
}

func (h Commands) logger(tx telebot.Context) zerolog.Logger {
	return tx.Get(loggerKey).(zerolog.Logger) //nolint:forcetypeassert
}

func (h Commands) menu(tx telebot.Context) *telebot.ReplyMarkup {
	menu, ok := tx.Get(menuKey).(*telebot.ReplyMarkup)

	if ok {
		return menu
	}

	return nil
}
