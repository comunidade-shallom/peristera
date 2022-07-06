package commands

import "gopkg.in/telebot.v3"

func (h Commands) Start(tx telebot.Context) error {
	return tx.Send(h.cfg.Description, telebot.ModeMarkdownV2)
}
