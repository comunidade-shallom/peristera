package config

import (
	"context"

	"github.com/comunidade-shallom/diakonos/pkg/sources"
)

type ctxKey struct{}

type AppConfig struct {
	Debug       bool            `fig:"-" yaml:"-"`
	Store       Store           `fig:"store" yaml:"store"`
	Timezone    string          `fig:"timezone" yaml:"timezone"`
	Description string          `fig:"description" yaml:"description"`
	Calendar    string          `fig:"calendar" yaml:"calendar"`
	Logger      Logger          `fig:"logger" yaml:"logger"`
	Telegram    Telegram        `fig:"telegram" yaml:"telegram"`
	Location    Location        `fig:"location" yaml:"location"`
	Youtube     YouTube         `fig:"youtube" yaml:"youtube"`
	Pix         Pix             `fig:"pix" yaml:"pix"`
	Cron        Cron            `fig:"cron" yaml:"cron"`
	Covers      sources.Sources `fig:"covers" yaml:"covers"`
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
