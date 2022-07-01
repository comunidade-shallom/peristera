package config

import (
	goErrors "errors"
	"os"
	"path"
	"path/filepath"

	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/creasty/defaults"
	"github.com/kkyr/fig"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

var (
	ErrFailToLoadConfig  = errors.System(nil, "fail to load config", "CONF:001")
	ErrFailEnsureConfig  = errors.System(nil, "fail to ensure config", "CONF:002")
	ConfigFileWasCreated = errors.Business("a new config file was created (%s)", "CONF:003")
)

func Load(file string) (AppConfig, error) {
	var err error

	cfg := AppConfig{}

	if file != "" {
		err = fig.Load(&cfg,
			fig.File(filepath.Base(file)),
			fig.Dirs(filepath.Dir(file)),
		)

		if err != nil {
			return cfg, ErrFailToLoadConfig.WithErr(err)
		}

		return applyDefaults(cfg)
	}

	home, err := homedir.Dir()
	if err != nil {
		return cfg, ErrFailToLoadConfig.WithErr(err)
	}

	err = fig.Load(&cfg,
		fig.File("peristera.yml"),
		fig.Dirs(
			".",
			path.Join(home, ".peristera"),
			path.Join(home, ".config"),
			path.Join(home, ".config/peristera"),
			home,
			"/etc/peristera",
			support.GetBinDirPath(),
		),
	)

	if goErrors.Is(err, fig.ErrFileNotFound) {
		return ensureConfig()
	}

	if err != nil {
		return cfg, err
	}

	return applyDefaults(cfg)
}

//nolint:cyclop
func applyDefaults(cfg AppConfig) (AppConfig, error) {
	// pwd, _ := os.Getwd()

	return cfg, nil
}

func ensureConfig() (AppConfig, error) {
	var err error

	cfg := AppConfig{}

	if err = defaults.Set(&cfg); err != nil {
		return cfg, ErrFailEnsureConfig.WithErr(err)
	}

	cfg, err = applyDefaults(cfg)

	if err != nil {
		return cfg, ErrFailEnsureConfig.WithErr(err)
	}

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		return cfg, ErrFailEnsureConfig.WithErr(err)
	}

	pwd, _ := os.Getwd()

	configFile := path.Join(pwd, "peristera.yml")

	err = os.WriteFile(configFile, buf, os.ModePerm)

	if err != nil {
		return cfg, ErrFailEnsureConfig.WithErr(err)
	}

	return cfg, ConfigFileWasCreated.Msgf(configFile)
}
