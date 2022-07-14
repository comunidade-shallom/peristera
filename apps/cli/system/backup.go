package system

import (
	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/system"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
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
			_ = system.NotifyBackupErr(cmd.Context, bot, cmd.Args().First(), roots, err)
		}

		db, err := database.Open(cfg.Store.Path)
		if err != nil {
			if force {
				notifyOnError(err)
			}

			return err
		}

		defer func() {
			logger.Warn().Err(db.Close()).Msg("Database closed.")
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

		return system.NotifyBackup(cmd.Context, bot, cmd.Args().First(), roots, file)
	},
}
