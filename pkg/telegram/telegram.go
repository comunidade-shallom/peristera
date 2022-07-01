package telegram

import (
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	tele "gopkg.in/telebot.v3"
)

const poolingTiming = 10 * time.Second

func NewBot(cfg config.AppConfig) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: poolingTiming},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return bot, err
	}

	bot.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	return bot, nil
}
