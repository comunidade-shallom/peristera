package telegram

import (
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"
)

func (h Handler) Me(tx telebot.Context) error {
	sender := tx.Sender()
	chat := tx.Chat()

	var builder strings.Builder

	builder.WriteString("*Name: *" + sender.FirstName + " " + sender.LastName)
	builder.WriteString("\n*Username: *" + sender.Username)
	builder.WriteString("\n*ID: * `" + strconv.Itoa(int(sender.ID)) + "`")
	builder.WriteString("\n\\-\\-\n")
	builder.WriteString("\n*Chat ID: * `" + strconv.Itoa(int(chat.ID)) + "`")

	return tx.Send(builder.String(), telebot.ModeMarkdownV2)
}
