package config

type Telegram struct {
	Token     string  `fig:"token" yaml:"token"`
	Broadcast []int64 `fig:"broadcast" yaml:"broadcast"`
}
