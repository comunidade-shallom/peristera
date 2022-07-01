package config

import (
	"context"
)

type ctxKey struct{}

type AppConfig struct {
	TelegramToken string `fig:"telegram_token" yaml:"telegram_token"`
}

func Ctx(ctx context.Context) AppConfig {
	cf, _ := ctx.Value(ctxKey{}).(AppConfig)

	return cf
}

func (c AppConfig) WithContext(ctx context.Context) context.Context {
	if cf, ok := ctx.Value(ctxKey{}).(AppConfig); ok {
		if cf == c {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, c)
}
