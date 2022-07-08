package system

import (
	"fmt"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/support"
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
		&cli.BoolFlag{
			Name:  "force",
			Usage: "force notify even if got error",
		},
	},
	Action: func(cmd *cli.Context) error {
		if !cmd.Bool("notify") {
			return ErrOnlyNotifyTrue
		}

		force := cmd.Bool("force")

		cfg := *config.Ctx(cmd.Context)

		roots := cfg.Telegram.Roots

		if len(roots) == 0 {
			return NoAdminsDefined
		}

		logger := zerolog.Ctx(cmd.Context).
			With().
			Str("context", "system").
			Logger()

		bot, err := telegram.NewBot(cfg)
		if err != nil {
			return err
		}

		notifyOnError := func(err error) {
			msg := fmt.Sprintf(
				"*%s*\n\n*System Notify:*\n`%s`\n\n*Backup Error:*\n`%s`\n\n*System time:*\n`%s`",
				support.AddSlashes(cmd.Args().First()),
				config.Hostname(),
				err.Error(),
				time.Now().Format(time.RFC3339),
			)

			for _, id := range roots {
				_, _ = bot.Send(&telebot.User{
					ID: id,
				}, msg, telebot.ModeMarkdownV2)
			}
		}

		db, err := database.Open(cfg.Store.Path)
		if err != nil {
			if force {
				notifyOnError(err)
			}

			return err
		}

		defer func() {
			_ = db.Close()
		}()

		logger.Debug().Msg("Generating backup...")

		file, destroy, err := system.Backup(cmd.Context, db)
		if err != nil {
			if force {
				notifyOnError(err)
			}

			return err
		}

		defer destroy()

		logger.Info().Msgf("Temp backup file created: %s", file.Name())

		st, err := file.Stat()
		if err != nil {
			if force {
				notifyOnError(err)
			}

			return err
		}

		if st.Size() == 0 {
			notifyOnError(BackupIsEmpty)

			return BackupIsEmpty
		}

		caption := fmt.Sprintf(
			"*%s*\n\n*System Notify:*\n`%s`\n\n*Peristera Backup:*\n`%s`",
			support.AddSlashes(cmd.Args().First()),
			config.Hostname(),
			time.Now().Format(time.RFC3339),
		)

		document := telegram.Document(file, caption)

		for _, id := range roots {
			document.Caption = caption // it changes over the loop, missing scapes

			_, err = bot.Send(&telebot.User{
				ID: id,
			}, document, telebot.ModeMarkdownV2)

			if err != nil {
				logger.Warn().Err(err).Int64("userId", id).Msg("Fail to sent backup file")
			}
		}

		logger.Info().Msgf("System backup sent do telegram %v", roots)

		return nil
	},
}
