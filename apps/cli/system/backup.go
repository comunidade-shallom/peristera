package system

import (
	"fmt"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/system"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var BackupCmd = &cli.Command{
	Name:  "backup",
	Usage: "do system backup",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "notify",
			Usage:    "send information to roots",
			Required: true,
		},
	},
	Action: func(cmd *cli.Context) error {
		if !cmd.Bool("notify") {
			return ErrOnlyNotifyTrue
		}

		cfg := *config.Ctx(cmd.Context)

		roots := cfg.Telegram.Roots

		if len(roots) == 0 {
			return NoAdminsDefined
		}

		logger := zerolog.Ctx(cmd.Context).
			With().
			Str("context", "system").
			Logger()

		db, err := database.Open(cfg.Store.Path)
		if err != nil {
			return err
		}

		defer func() {
			_ = db.Close()
		}()

		logger.Debug().Msg("Generating backup...")

		file, destroy, err := system.Backup(cmd.Context, db)
		if err != nil {
			return err
		}

		defer destroy()

		logger.Info().Msgf("Temp backup file created: %s", file.Name())

		bot, err := telegram.NewBot(cfg)
		if err != nil {
			return err
		}

		document := telegram.Document(
			file,
			fmt.Sprintf("*System Notify*\n\n*Peristera Backup:\n*%s", time.Now().Format(time.RFC3339)),
		)

		for _, id := range roots {
			_, err = bot.Send(&telebot.User{
				ID: id,
			}, document, telebot.ModeMarkdownV2)

			if err != nil {
				return err
			}
		}

		logger.Info().Msgf("System backup sent do telegram %v", roots)

		return nil
	},
}
