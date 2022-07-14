package commands

import (
	"strconv"

	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

const (
	loggerKey = "logger"
	menuKey   = "menu"
)

func useLogger(parent zerolog.Logger) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(tx telebot.Context) error {
			logger := parent.With().
				Str("message_id", strconv.Itoa(tx.Message().ID)).
				Str("sender_id", strconv.FormatInt(tx.Sender().ID, 10)). //nolint:gomnd
				Str("chat_id", strconv.FormatInt(tx.Chat().ID, 10)).     //nolint:gomnd
				Logger()

			tx.Set(loggerKey, logger)

			cmd := tx.Message().Text

			logger.Info().Msgf("New trigger (%s)", cmd)

			return next(tx)
		}
	}
}

func useMenu(menu *telebot.ReplyMarkup) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(tx telebot.Context) error {
			tx.Set(menuKey, menu)

			return next(tx)
		}
	}
}

func restrictTo(ids []int64, role string) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(tx telebot.Context) error {
			senderID := tx.Sender().ID

			for _, ID := range ids {
				if ID == senderID {
					return next(tx)
				}
			}

			logger := tx.Get(loggerKey).(zerolog.Logger) //nolint:forcetypeassert

			logger.Warn().Msg("Rejecting command from non " + role)

			err := tx.Reply("This command is only for " + role)
			if err != nil {
				logger.Error().Err(err).Msg("Fail to reply")

				return err
			}

			return nil
		}
	}
}
