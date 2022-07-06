package telegram

import (
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"gopkg.in/telebot.v3"
)

const poolingTiming = 10 * time.Second

func NewBot(cfg config.AppConfig) (*telebot.Bot, error) {
	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: poolingTiming},
	}

	return telebot.NewBot(pref)
}
