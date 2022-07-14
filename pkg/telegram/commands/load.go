package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/comunidade-shallom/peristera/pkg/support"
	"gopkg.in/telebot.v3"
)

func (h Commands) Load(tx telebot.Context) error {
	if !tx.Message().IsReply() {
		return tx.Reply("You must replay a backup message.")
	}

	document := tx.Message().ReplyTo.Document

	if document == nil {
		return tx.Reply("Document is not defined.")
	}

	if !strings.HasPrefix(document.FileName, "peristera.") {
		return tx.Reply(
			fmt.Sprintf("Incompatible document: %s", support.AddSlashes(document.FileName)),
			telebot.ModeMarkdownV2,
		)
	}

	if err := tx.Notify(telebot.Typing); err != nil {
		return err
	}

	file, err := ioutil.TempFile(os.TempDir(), "peristera.restore.*.bak")
	if err != nil {
		return err
	}

	tempFileName := file.Name()

	defer func() {
		_ = file.Close()
		_ = os.Remove(tempFileName)
	}()

	logger := h.logger(tx)

	if err = tx.Bot().Download(&document.File, file.Name()); err != nil {
		return err
	}

	logger.Info().Msgf("Temp file generated: %s", file.Name())

	file, err = os.Open(tempFileName)
	if err != nil {
		return err
	}

	if err = h.db.Load(file); err != nil {
		return err
	}

	logger.Info().Msgf("Database loaded from: %s", file.Name())

	return tx.Reply("Database loaded with backup.")
}
