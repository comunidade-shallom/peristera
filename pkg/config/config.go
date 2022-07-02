package config

import (
	"context"
)

type ctxKey struct{}

type AppConfig struct {
	Telegram Telegram `fig:"token" yaml:"token"`
	Youtube  YouTube  `fig:"youtube" yaml:"youtube"`
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
