package commands

import (
	"image/png"
	"os"
	"strings"

	"github.com/comunidade-shallom/diakonos/pkg/covers"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"gopkg.in/telebot.v3"
)

var ErrMissingText = errors.Business("Text must be defined", "TC:002")

func (h Commands) Cover(tx telebot.Context) error {
	logger := h.logger(tx)

	payload := strings.TrimSpace(tx.Message().Payload)

	if len(payload) == 0 {
		return ErrMissingText
	}

	if err := tx.Reply("Generating cover..."); err != nil {
		return err
	}

	if err := tx.Notify(telebot.Typing); err != nil {
		return err
	}

	logger.Info().Msgf("Generating cover: %s", payload)

	cover, err := covers.GeneratorSource{
		Sources: h.cfg.Covers,
		Text:    payload,
	}.Random()

	if err != nil {
		logger.Error().Err(err).Msg("Fail to generate cover")

		return err
	}

	file, err := os.CreateTemp(os.TempDir(), "peristera-cover-*.png")

	if err != nil {
		logger.Error().Err(err).Msg("Fail to generate tmp dir")

		return err
	} else {
		logger.Debug().Msgf("Temp file created: %s", file.Name())
	}

	defer file.Close()
	defer os.Remove(file.Name())

	png.Encode(file, cover.Build())

	if err := file.Sync(); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	logger.Info().Msgf("Cover generated: %s", file.Name())

	photo := &telebot.Photo{
		File: telebot.FromReader(file),
	}

	return tx.SendAlbum(telebot.Album{photo}, telebot.ModeMarkdownV2)
}
