package commands

import "gopkg.in/telebot.v3"

func (h Commands) Start(tx telebot.Context) error {
	return tx.Reply(
		h.cfg.Description,
		h.menu(tx),
		telebot.ModeMarkdownV2,
	)
}
