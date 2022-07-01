package telegram

import (
	"strconv"
	"strings"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	cfg config.AppConfig
}

func (h Handler) Me(ctx telebot.Context) error {
	sender := ctx.Sender()

	var builder strings.Builder

	builder.WriteString("*Name: *" + sender.FirstName + " " + sender.LastName)
	builder.WriteString("\n*Username: *" + sender.Username)
	builder.WriteString("\n*ID: * `" + strconv.Itoa(int(sender.ID)) + "`")

	return ctx.Send(builder.String(), telebot.ModeMarkdownV2)
}
