package commands

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
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

	bot.Use(middleware.Recover())
	bot.Use(useLogger(logger))

	_, err := h.registerMenu(ctx, bot)
	if err != nil {
		return err
	}

	err = h.registerCommands(ctx, bot)
	if err != nil {
		return err
	}

	adm := bot.Group()
	adm.Use(restrictTo(h.cfg.Telegram.Admins, "admins"))
	adm.Handle("/me", h.Me)
	adm.Handle("/system", h.System)
	adm.Handle("/cover", h.Cover)

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

func (h Commands) registerMenu(ctx context.Context, bot *telebot.Bot) (*telebot.ReplyMarkup, error) {
	logger := zerolog.Ctx(ctx).With().Str("fn", "commands:registerMenu").Logger()

	menu := &telebot.ReplyMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
	}

	bot.Use(useMenu(menu))

	items := h.cfg.Telegram.Commands.Menu

	buttons := make([]telebot.Btn, len(items))

	for index, item := range items {
		handler, err := h.getHandler(item.Handler)
		if err != nil {
			return menu, err
		}

		btn := menu.Text(item.Text)

		bot.Handle(&btn, handler)

		buttons[index] = btn

		logger.Info().Msgf("Handler '%s' registred to menu '%s'", item.Handler, item.Text)
	}

	menu.Reply(menu.Split(2, buttons)...) //nolint:gomnd

	return menu, nil
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
	return tx.Get(telegram.LoggerKey).(zerolog.Logger) //nolint:forcetypeassert
}

func (h Commands) menu(tx telebot.Context) *telebot.ReplyMarkup {
	menu, ok := tx.Get(menuKey).(*telebot.ReplyMarkup)

	if ok {
		return menu
	}

	return nil
}
