package config

type Logger struct {
	Level  string `fig:"level" yaml:"level" default:"info"`
	Format string `fig:"format" yaml:"format" default:"text"`
}

func (l Logger) Debug(level string) bool {
	debugLevels := [2]string{"debug", "trace"}

	for _, val := range debugLevels {
		if val == level {
			return true
		}
	}

	for _, val := range debugLevels {
		if val == l.Level {
			return true
		}
	}

	return false
}
