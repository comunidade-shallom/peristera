package config

import (
	"context"
)

type ctxKey struct{}

type Channel struct {
	Name string `fig:"name" yaml:"name"`
	ID   string `fig:"id" yaml:"id"`
}

type AppConfig struct {
	TelegramToken string    `fig:"telegram_token" yaml:"telegram_token"`
	YoutubeToken  string    `fig:"youtube_token" yaml:"youtube_token"`
	Channels      []Channel `fig:"channels" yaml:"channels"`
}

func Ctx(ctx context.Context) *AppConfig {
	cf, _ := ctx.Value(ctxKey{}).(*AppConfig)

	return cf
}

func (c *AppConfig) WithContext(ctx context.Context) context.Context {
	if cf, ok := ctx.Value(ctxKey{}).(*AppConfig); ok {
		if cf == c {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, c)
}
