package telegram

import (
	"context"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	tele "gopkg.in/telebot.v3"
)

const poolingTiming = 10 * time.Second

func NewBot(ctx context.Context, cfg config.AppConfig) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: poolingTiming},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return bot, err
	}

	youtube, err := ytube.NewService(ctx, cfg)
	if err != nil {
		return bot, err
	}

	handlers := Handler{
		cfg:     cfg,
		youtube: youtube,
	}

	bot.Handle("/me", handlers.Me)
	bot.Handle("/videos", handlers.Videos)

	return bot, nil
}
