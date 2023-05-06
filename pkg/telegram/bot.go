package telegram

import (
	"fmt"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

const poolingTiming = 10 * time.Second

const LoggerKey = "bot:logger"

func NewBot(cfg config.AppConfig) (*telebot.Bot, error) {
	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: poolingTiming},
		OnError: func(err error, tx telebot.Context) {
			_ = tx.Reply(fmt.Sprintf("Error: %s", err.Error()))
			logger := tx.Get(LoggerKey).(zerolog.Logger) //nolint:forcetypeassert
			logger.Error().Err(err).Msg("Bot error")
		},
	}

	return telebot.NewBot(pref)
}
