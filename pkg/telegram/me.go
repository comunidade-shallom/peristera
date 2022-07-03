package telegram

import (
	"strconv"
	"strings"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	logger  zerolog.Logger
	cfg     config.AppConfig
	youtube ytube.Service
	bot     *telebot.Bot
}

func (h Handler) Me(ctx telebot.Context) error {
	sender := ctx.Sender()
	chat := ctx.Chat()

	var builder strings.Builder

	builder.WriteString("*Name: *" + sender.FirstName + " " + sender.LastName)
	builder.WriteString("\n*Username: *" + sender.Username)
	builder.WriteString("\n*ID: * `" + strconv.Itoa(int(sender.ID)) + "`")
	builder.WriteString("\n\\-\\-\n")
	builder.WriteString("\n*Chat ID: * `" + strconv.Itoa(int(chat.ID)) + "`")

	return ctx.Send(builder.String(), telebot.ModeMarkdownV2)
}
