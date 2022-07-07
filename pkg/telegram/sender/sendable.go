package sender

import (
	"context"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/database"
	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type Recipients interface {
	Recipients() []telebot.Recipient
}

type Message interface {
	Hash() []byte
	Params() []interface{}
	Content() interface{}
}

type Sendable interface {
	Recipients
	Message
}

func SendableWorker(ctx context.Context, in <-chan Sendable, bot *telebot.Bot, db database.Database) {
	parent := zerolog.Ctx(ctx).With().Str("worker", "bot-sender").Logger()
	parent.Info().Msg("Starting sendable worker")

	for msg := range in {
		current := msg
		logger := parent.With().Bytes("hash", current.Hash()).Logger()

		logger.Info().Msg("Sending message...")

		for _, recipient := range current.Recipients() {
			_, err := bot.Send(recipient, current.Content(), current.Params()...)
			if err != nil {
				logger.Error().
					Err(err).
					Str("Recipient", recipient.Recipient()).
					Msg("Fail to send message")
			} else {
				logger.Debug().
					Str("Recipient", recipient.Recipient()).
					Msg("Message sent")
			}
		}

		err := db.DB().Update(func(txn *badger.Txn) error {
			return txn.Set(current.Hash(), []byte(time.Now().Format(time.RFC3339)))
		})
		if err != nil {
			logger.Error().
				Err(err).
				Bytes("hash", current.Hash()).
				Msg("Fail store message")
		} else {
			logger.Info().Msg("Message stored...")
		}
	}

	parent.Warn().Msg("Sendable worker is stopped")
}
