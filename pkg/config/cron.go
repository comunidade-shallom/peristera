package config

type Cron struct {
	Backup      []string `fig:"backup" yaml:"backup"`
	LastUpdates []string `fig:"last_updates" yaml:"last_updates"`
}
