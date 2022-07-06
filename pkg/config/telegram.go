package config

type Telegram struct {
	Token     string  `fig:"token" yaml:"token"`
	Broadcast []int64 `fig:"broadcast" yaml:"broadcast"`
	Admins    []int64 `fig:"admins" yaml:"admins"`
	Roots     []int64 `fig:"roots" yaml:"roots"`
}
