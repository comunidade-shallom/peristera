package cron

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"gopkg.in/telebot.v3"
)

func New(ctx context.Context, cfg config.AppConfig, bot *telebot.Bot, youtube ytube.Service) (Jobs, error) {
	jb := Jobs{
		youtube: youtube,
		cfg:     cfg,
		bot:     bot,
	}

	return jb, jb.register(ctx)
}
