package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
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

	defer destroy()

	logger.Info().Msgf("Temp backup file created: %s", file.Name())

	st, err := file.Stat()
	if err != nil {
		return err
	}

	if st.Size() == 0 {
		logger.Warn().Msg("Backup file is empty")

		return tx.Reply(
			fmt.Sprintf("🪣 Backup file is empty\n\n📁 `%s`", h.cfg.Store.Path),
			telebot.ModeMarkdownV2,
		)
	}

	logger.Debug().Msg("Backup write on disk")

	if err = tx.Notify(telebot.UploadingDocument); err != nil {
		return err
	}

	caption := fmt.Sprintf(
		"*🖥️ System:* `%s`\n\n *🗄️ Peristera Backup:*\n `%s`\n\n📁 `%s`",
		config.Hostname(),
		time.Now().Format(time.RFC3339),
		h.cfg.Store.Path,
	)

	document := telegram.Document(file, caption)

	if err = tx.Reply(document, telebot.ModeMarkdownV2); err != nil {
		return err
	}

	logger.Info().Msg("Backup done")

	return nil
}
