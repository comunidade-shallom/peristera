package config

import (
	goErrors "errors"
	"os"
	"path"
	"path/filepath"

	"github.com/comunidade-shallom/diakonos/pkg/sources"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/creasty/defaults"
	"github.com/kkyr/fig"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

var (
	ErrFailToLoadConfig     = errors.System(nil, "fail to load config", "CONF:001")
	ErrFailEnsureConfig     = errors.System(nil, "fail to ensure config", "CONF:002")
	ConfigFileWasCreated    = errors.Business("a new config file was created (%s)", "CONF:003")
	ErrMissingTelegramToken = errors.System(nil, "missing telegram token", "CONF:004")
	ErrMissingYoutubeToken  = errors.System(nil, "missing youtube token", "CONF:005")
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
			"/peristera.d",
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

func applyDefaults(cfg AppConfig) (AppConfig, error) {
	if cfg.Telegram.Token == "" {
		cfg.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	}

	if cfg.Youtube.Token == "" {
		cfg.Youtube.Token = os.Getenv("YOUTUBE_TOKEN")
	}

	if cfg.Telegram.Token == "" {
		return cfg, ErrMissingTelegramToken
	}

	if cfg.Youtube.Token == "" {
		return cfg, ErrMissingYoutubeToken
	}

	if len(cfg.Covers.Colors) == 0 {
		cfg.Covers.Colors = sources.DefaultColors()
	}

	wd, err := os.Getwd()
	if err != nil {
		return cfg, err
	}

	if !filepath.IsAbs(cfg.Covers.Covers) {
		cfg.Covers.Covers = filepath.Join(wd, cfg.Covers.Covers)
	}

	if !filepath.IsAbs(cfg.Covers.Fonts) {
		cfg.Covers.Fonts = filepath.Join(wd, cfg.Covers.Fonts)
	}

	if !filepath.IsAbs(cfg.Covers.Footer) {
		cfg.Covers.Footer = filepath.Join(wd, cfg.Covers.Footer)
	}

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
