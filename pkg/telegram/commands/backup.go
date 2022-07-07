package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/system"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"gopkg.in/telebot.v3"
)

func (h Commands) Backup(tx telebot.Context) error {
	logger := h.logger(tx)

	if err := tx.Reply("Generating backup..."); err != nil {
		return err
	}

	if err := tx.Notify(telebot.Typing); err != nil {
		return err
	}

	file, destroy, err := system.Backup(logger.WithContext(context.TODO()), h.db)
	if err != nil {
		return err
	}

	logger.Info().Msgf("Temp backup file created: %s", file.Name())

	defer destroy()

	logger.Debug().Msg("Backup write on disk")

	if err = tx.Notify(telebot.UploadingDocument); err != nil {
		return err
	}

	document := telegram.Document(file, fmt.Sprintf("*Peristera Backup:*%s", time.Now().Format(time.RFC3339)))

	if err = tx.Reply(document, telebot.ModeMarkdownV2); err != nil {
		return err
	}

	logger.Info().Msg("Backup done")

	return nil
}
