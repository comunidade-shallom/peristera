package sender

import (
	"context"

	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type Recipients interface {
	Recipients() []telebot.Recipient
}

type Message interface {
	Hash() string
	Params() []interface{}
	Content() interface{}
}

type Sendable interface {
	Recipients
	Message
}

func SendableWorker(ctx context.Context, in <-chan Sendable, bot *telebot.Bot) {
	logger := zerolog.Ctx(ctx).With().Str("worker", "bot-sender").Logger()
	logger.Info().Msg("Starting sendable worker")

	for msg := range in {
		logger.Info().Str("hash", msg.Hash()).Msg("Sending message...")

		for _, recipient := range msg.Recipients() {
			_, err := bot.Send(recipient, msg.Content(), msg.Params()...)
			if err != nil {
				logger.Error().Err(err).Msg("Fail to send message")
			}
		}
	}

	logger.Warn().Msg("Sendable worker is stopped")
}
