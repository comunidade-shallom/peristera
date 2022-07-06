package commands

import (
	"strconv"
	"strings"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

const loggerKey = "logger"

func useLogger(parent zerolog.Logger) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(tx telebot.Context) error {
			logger := parent.With().
				Str("message_id", strconv.Itoa(tx.Message().ID)).
				Str("sender_id", strconv.FormatInt(tx.Sender().ID, 10)). //nolint:gomnd
				Str("chat_id", strconv.FormatInt(tx.Chat().ID, 10)).     //nolint:gomnd
				Logger()

			tx.Set(loggerKey, logger)

			cmd := "--"

			if text := tx.Message().Text; strings.HasPrefix(text, "/") {
				cmd = text
			}

			logger.Info().Msgf("New trigger (%s)", cmd)

			return next(tx)
		}
	}
}

func onlyAdmins(cfg config.Telegram) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(tx telebot.Context) error {
			senderID := tx.Sender().ID

			for _, ID := range cfg.Admins {
				if ID == senderID {
					return next(tx)
				}
			}

			logger := tx.Get(loggerKey).(zerolog.Logger) //nolint:forcetypeassert

			logger.Warn().Msg("Rejecting command from non admin")

			err := tx.Reply("This command is only for admins")
			if err != nil {
				logger.Error().Err(err).Msg("Fail to reply")
			}

			return nil
		}
	}
}
