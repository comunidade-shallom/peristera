package system

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type DestroyBackup func()

func noop() {}

func Backup(ctx context.Context, db database.Database) (*os.File, DestroyBackup, error) {
	logger := zerolog.Ctx(ctx)

	file, err := os.CreateTemp(os.TempDir(), "peristera.*.bak")
	if err != nil {
		return file, noop, err
	}

	destroy := func() {
		name := file.Name()
		_ = file.Close()
		_ = os.Remove(name)
	}

	logger.Info().Msgf("Temp backup file created: %s", file.Name())

	bw := bufio.NewWriterSize(file, 64<<20) //nolint:gomnd

	// run backup
	if err = db.Backup(bw); err != nil {
		return file, destroy, err
	}

	logger.Debug().Msg("Backup generated")

	if err = bw.Flush(); err != nil {
		return file, destroy, err
	}

	if err = file.Sync(); err != nil {
		return file, destroy, err
	}

	return file, destroy, nil
}

func NotifyBackupErr(ctx context.Context, bot *telebot.Bot, header string, userIDs []int64, err error) error {
	host := config.Hostname() + " (" + config.Version() + ")"
	cfg := *config.Ctx(ctx)

	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "system:notify-backup-err").
		Logger()

	logger.Warn().Msg("Notify on error...")

	msg := fmt.Sprintf(
		"*%s*\n\n*🖥️ System Notify:*\n`%s`\n\n*🔴 Backup Error:*\n`%s`\n\n*🕰️ System time:*\n`%s`\n\n📁 `%s`",
		support.AddSlashes(header),
		host,
		err.Error(),
		time.Now().Format(time.RFC3339),
		cfg.Store.Path,
	)

	for _, id := range userIDs {
		_, err = bot.Send(&telebot.User{
			ID: id,
		}, msg, telebot.ModeMarkdownV2)

		if err != nil {
			logger.Warn().Err(err).Int64("userId", id).Msg("Fail to sent notify")
		}
	}

	return nil
}

func NotifyBackup(ctx context.Context, bot *telebot.Bot, header string, userIDs []int64, file *os.File) error {
	host := config.Hostname() + " (" + config.Version() + ")"

	cfg := *config.Ctx(ctx)

	logger := zerolog.Ctx(ctx).
		With().
		Str("context", "system:notify-backup").
		Logger()

	caption := fmt.Sprintf(
		"*%s*\n\n*🖥️ System Notify:*\n`%s`\n\n*🗄️ Peristera Backup:*\n`%s`\n\n📁`%s`",
		support.AddSlashes(header),
		host,
		time.Now().Format(time.RFC3339),
		cfg.Store.Path,
	)

	document := telegram.Document(file, caption)

	for _, id := range userIDs {
		document.Caption = caption // it changes over the loop, missing scapes

		_, err := bot.Send(&telebot.User{
			ID: id,
		}, document, telebot.ModeMarkdownV2)
		if err != nil {
			logger.Warn().Err(err).Int64("userId", id).Msg("Fail to sent backup file")
		}
	}

	logger.Info().Msgf("System backup sent do telegram %v", userIDs)

	return nil
}
