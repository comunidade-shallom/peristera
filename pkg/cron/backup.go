package cron

import (
	"context"

	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/comunidade-shallom/peristera/pkg/system"
)

var (
	NoRootsDefined = errors.Business("No roots defined", "CR:002")
	BackupIsEmpty  = errors.Business("🪣 Backup is empty", "CR:003")
)

func (j *Jobs) Backup(ctx context.Context) error {
	roots := j.cfg.Telegram.Roots

	if len(roots) == 0 {
		return NoRootsDefined
	}

	header := "Backup scheduler"

	notifyOnError := func(err error) {
		_ = system.NotifyBackupErr(ctx, j.bot, header, roots, err)
	}

	logger := j.jobLogger(ctx, "backup")

	logger.Debug().Msg("Generating backup...")

	file, destroy, err := system.Backup(ctx, j.database)
	if err != nil {
		notifyOnError(err)

		return err
	}

	defer destroy()

	logger.Info().Msgf("Temp backup file created: %s", file.Name())

	st, err := file.Stat()
	if err != nil {
		notifyOnError(err)

		return err
	}

	if st.Size() == 0 {
		notifyOnError(BackupIsEmpty)

		return BackupIsEmpty
	}

	return system.NotifyBackup(ctx, j.bot, header, roots, file)
}
