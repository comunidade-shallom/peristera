package config

import (
	"context"
)

type ctxKey struct{}

type AppConfig struct {
	Debug       bool     `fig:"-" yaml:"-"`
	Logger      Logger   `fig:"logger" yaml:"logger"`
	Telegram    Telegram `fig:"token" yaml:"token"`
	Youtube     YouTube  `fig:"youtube" yaml:"youtube"`
	Pix         Pix      `fig:"pix" yaml:"pix"`
	Description string   `fig:"description" yaml:"description"`
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

func (c AppConfig) Tags() map[string]interface{} {
	return map[string]interface{}{}
}
