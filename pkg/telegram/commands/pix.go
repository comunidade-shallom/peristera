package commands

import (
	"gopkg.in/telebot.v3"
)

func (h Commands) Pix(tx telebot.Context) error {
	pix := h.cfg.Pix

	if len(pix.QRCode) == 0 {
		return tx.Reply(pix.Description, telebot.ModeMarkdownV2)
	}

	buf, err := pix.QRCode.NewBuffer()
	if err != nil {
		return err
	}

	photo := &telebot.Photo{
		File:    telebot.FromReader(buf),
		Caption: pix.Description,
	}

	return tx.SendAlbum(telebot.Album{photo}, telebot.ModeMarkdownV2)
}
