package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

const timeout = 5
const maxResults = 2

type Handler struct {
	logger  zerolog.Logger
	cfg     config.AppConfig
	youtube ytube.Service
}

func (h Handler) Me(ctx telebot.Context) error {
	sender := ctx.Sender()

	var builder strings.Builder

	builder.WriteString("*Name: *" + sender.FirstName + " " + sender.LastName)
	builder.WriteString("\n*Username: *" + sender.Username)
	builder.WriteString("\n*ID: * `" + strconv.Itoa(int(sender.ID)) + "`")

	return ctx.Send(builder.String(), telebot.ModeMarkdownV2)
}

func (h Handler) Videos(ctx telebot.Context) error {
	_ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)

	logger := ctx.Get(loggerKey).(zerolog.Logger)

	defer cancel()

	for _, ch := range h.cfg.Youtube.Channels {
		if err := ctx.Notify(telebot.Typing); err != nil {
			return err
		}

		logger.Info().Msgf("Loading last videos %s", ch.Name)

		vids, err := h.youtube.LastVideos(_ctx, ch.ID, maxResults)
		if err != nil {
			return err
		}

		if len(vids) == 0 {
			msg := fmt.Sprintf("%s: Sem resultados", ch.Name)
			logger.Warn().Msg(msg)

			if err = ctx.Reply(msg); err != nil {
				return err
			}

			continue
		}

		for _, vid := range vids {
			ctx.Send(vid.URL())
		}

		logger.Info().Msgf("%v Videos from %s", len(vids), ch.Name)
	}

	return nil
}
