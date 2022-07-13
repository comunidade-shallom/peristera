package commands

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

var ErrHandlerNotFound = errors.Business("handler %s not found", "TC:001")

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

	if err := h.registerCommands(ctx, bot); err != nil {
		return err
	}

	adm := bot.Group()
	adm.Use(restrictTo(h.cfg.Telegram.Admins, "admins"))

	adm.Handle("/me", h.Me)
	adm.Handle("/system", h.System)

	root := bot.Group()
	root.Use(restrictTo(h.cfg.Telegram.Roots, "roots"))
	root.Handle("/exec", h.Exec)
	root.Handle("/backup", h.Backup)
	root.Handle("/load", h.Load)

	return nil
}

func (h Commands) getHandler(name string) (telebot.HandlerFunc, error) {
	switch name {
	case "start":
		return h.Start, nil
	case "pix":
		return h.Pix, nil
	case "videos":
		return h.Videos, nil
	case "address":
		return h.Address, nil
	case "calendar":
		return h.Calendar, nil
	default:
		return nil, ErrHandlerNotFound.Msgf(name)
	}
}

func (h Commands) registerMenu(bot *telebot.Bot) *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}

	bot.Use(useMenu(menu))

	btnAbout := menu.Text("ℹ️ Sobre")
	btnAgenda := menu.Text("🗓️ Agenda")
	btnAddress := menu.Text("📍 Endereço")
	btnPix := menu.Text("🏦 Pix")
	btnYoutube := menu.Text("📺 YouTube")

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

func (h Commands) registerCommands(ctx context.Context, bot *telebot.Bot) error {
	logger := zerolog.Ctx(ctx).With().Str("fn", "commands:registerCommands").Logger()

	cfg := h.cfg.Telegram
	for _, mapper := range cfg.Commands.Mappers {
		handler, err := h.getHandler(mapper.Handler)
		if err != nil {
			return err
		}

		for _, endpoint := range mapper.Endpoints {
			bot.Handle(endpoint, handler)
		}

		logger.Info().Msgf("Handler '%s' registred to %v", mapper.Handler, mapper.Endpoints)
	}

	cmds := make([]telebot.Command, len(cfg.Commands.SetOf))

	for index, set := range cfg.Commands.SetOf {
		cmds[index] = telebot.Command{
			Text:        set.Text,
			Description: set.Description,
		}
	}

	return bot.SetCommands(cmds)
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
