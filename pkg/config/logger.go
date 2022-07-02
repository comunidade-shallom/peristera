package config

type Logger struct {
	Level  string `fig:"level" yaml:"level" default:"debug"`
	Format string `fig:"format" yaml:"format" default:"text"`
}
