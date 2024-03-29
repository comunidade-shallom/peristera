//nolint:gofumpt
package commands

import (
	"image/png"
	"os"
	"regexp"
	"strings"

	"github.com/comunidade-shallom/diakonos/pkg/covers"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"gopkg.in/telebot.v3"
)

var ErrMissingText = errors.Business("Text must be defined", "TC:002")

var sizeRegx = regexp.MustCompile(`^\d+x\d+`)

const defaultSize = 1080

func (h Commands) Cover(tx telebot.Context) error {
	logger := h.logger(tx)

	size, text := BuildCoverParams(tx.Message())

	if len(text) == 0 {
		return ErrMissingText
	}

	if err := tx.Reply("Generating cover image " + size.String() + "..."); err != nil {
		return err
	}

	if err := tx.Notify(telebot.Typing); err != nil {
		return err
	}

	logger.Info().Msgf("Generating cover image: %s", text)

	return h.buildCover(tx, size, text)
}

func (h Commands) buildCover(tx telebot.Context, size covers.Size, text string) error {
	logger := h.logger(tx)

	cover, err := covers.GeneratorSource{
		Sources: h.cfg.Covers,
		Width:   size.Width,
		Height:  size.Height,
		Text:    text,
	}.Random()

	if err != nil {
		logger.Error().Err(err).Msg("Fail to generate cover")

		return err
	}

	file, err := os.CreateTemp(os.TempDir(), "peristera-cover-*.png")

	if err != nil {
		logger.Error().Err(err).Msg("Fail to generate tmp dir")

		return err
	}

	logger.Debug().Msgf("Temp file created: %s", file.Name())

	defer file.Close()
	defer os.Remove(file.Name())

	if err = png.Encode(file, cover.Build()); err != nil {
		logger.Error().Err(err).Msg("Fail to encode cover")

		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	logger.Info().Msgf("Cover generated: %s", file.Name())

	photo := &telebot.Photo{
		File:    telebot.FromReader(file),
		Caption: text,
	}

	if err := tx.Notify(telebot.UploadingPhoto); err != nil {
		return err
	}

	return tx.SendAlbum(telebot.Album{photo})
}

func parseCoverPayload(raw string) (covers.Size, string) {
	matchs := sizeRegx.FindAllString(raw, 1)
	size := covers.Size{
		Width:  defaultSize,
		Height: defaultSize,
	}

	if len(matchs) > 0 {
		size = covers.ParseSize(matchs[0])
		raw = sizeRegx.ReplaceAllString(raw, "")
	}

	return size, strings.TrimSpace(raw)
}

func BuildCoverParams(msg *telebot.Message) (covers.Size, string) {
	var text string

	var size covers.Size

	payload := strings.TrimSpace(msg.Payload)

	if replyTo := msg.ReplyTo; replyTo == nil {
		size, text = parseCoverPayload(payload)
	} else if replyTo.Photo != nil {
		text = replyTo.Caption
		size = covers.Size{
			Width:  replyTo.Photo.Width,
			Height: replyTo.Photo.Height,
		}
	} else {
		text = replyTo.Text
		size = covers.ParseSize(payload)
	}

	if size.Width == 0 {
		size.Width = defaultSize
	}

	if size.Height == 0 {
		size.Height = size.Width
	}

	return size, text
}
