package telegram

import "gopkg.in/telebot.v3"

func (h Handler) Start(tx telebot.Context) error {
	return tx.Send(h.cfg.Description, telebot.ModeMarkdownV2)
}
