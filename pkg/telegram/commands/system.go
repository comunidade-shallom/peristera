package commands

import (
	"github.com/comunidade-shallom/peristera/pkg/support/system"
	"gopkg.in/telebot.v3"
)

func (h Commands) System(tx telebot.Context) error {
	logger := h.logger(tx)

	logger.Info().Msg("Loading system info...")

	data, err := system.New()
	if err != nil {
		return err
	}

	return tx.Reply(data.MarkdownV2(""), telebot.ModeMarkdownV2)
}
