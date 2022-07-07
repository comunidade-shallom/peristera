package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/support"
	"gopkg.in/telebot.v3"
)

func (h Commands) Backup(tx telebot.Context) error {
	file, err := ioutil.TempFile(os.TempDir(), "peristera.*.bak")
	if err != nil {
		return err
	}

	logger := h.logger(tx)

	logger.Info().Msgf("Temp backup file created: %s", file.Name())

	if err = tx.Reply("Generating backup..."); err != nil {
		return err
	}

	if err = tx.Notify(telebot.Typing); err != nil {
		return err
	}

	bw := bufio.NewWriterSize(file, 64<<20) //nolint:gomnd

	// run backup
	if err = h.db.Backup(bw); err != nil {
		return err
	}

	logger.Debug().Msg("Backup generated")

	if err = bw.Flush(); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	defer func() {
		name := file.Name()
		_ = file.Close()
		_ = os.Remove(name)
	}()

	logger.Debug().Msg("Backup write on disk")

	if err = tx.Notify(telebot.UploadingDocument); err != nil {
		return err
	}

	document := &telebot.Document{File: telebot.FromDisk(file.Name())}
	document.Caption = fmt.Sprintf("*Peristera Backup:* %s", support.AddSlashes(time.Now().Format(time.RFC3339)))
	document.FileName = file.Name()

	if err = tx.Reply(document, telebot.ModeMarkdownV2); err != nil {
		return err
	}

	logger.Info().Msg("Backup done")

	return nil
}
