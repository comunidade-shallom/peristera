package telegram

import (
	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	cfg     config.AppConfig
	youtube ytube.Service
	bot     *telebot.Bot
}

func (h Handler) logger(tx telebot.Context) zerolog.Logger {
	return tx.Get(loggerKey).(zerolog.Logger) //nolint:forcetypeassert
}
