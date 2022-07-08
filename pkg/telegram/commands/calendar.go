package commands

import "gopkg.in/telebot.v3"

func (h Commands) Calendar(tx telebot.Context) error {
	return tx.Reply(h.cfg.Calendar, telebot.ModeMarkdownV2)
}
