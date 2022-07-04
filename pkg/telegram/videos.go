package telegram

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

const (
	loadTimeout         = 5
	defaultCountResults = 2
	maxResults          = 5
)

func (h Handler) Videos(tx telebot.Context) error {
	_ctx, cancel := context.WithTimeout(context.Background(), time.Second*loadTimeout)

	logger := tx.Get(loggerKey).(zerolog.Logger) //nolint:forcetypeassert

	defer cancel()

	count := getMaxArg(tx.Args(), logger, tx)

	for _, ch := range h.cfg.Youtube.Channels {
		if err := tx.Notify(telebot.Typing); err != nil {
			return err
		}

		logger.Info().Msgf("Loading last videos %s", ch.Name)

		vids, err := h.youtube.LastVideos(_ctx, ch.ID, count)
		if err != nil {
			return err
		}

		if len(vids) == 0 {
			msg := fmt.Sprintf("%s: Sem resultados", ch.Name)
			logger.Warn().Msg(msg)

			if err = tx.Reply(msg); err != nil {
				return err
			}

			continue
		}

		for _, vid := range vids {
			err = tx.Send(vid.URL())

			if err != nil {
				logger.Warn().Err(err).Msg("Error on send")
			}
		}

		err = tx.Send(
			fmt.Sprintf("Mais vídeos de %s\n\n%s", ch.Name, ch.GetURL()),
		)

		if err != nil {
			logger.Warn().Err(err).Msg("Error on send btn")
		}

		logger.Info().Msgf("%v Videos from %s", len(vids), ch.Name)
	}

	return nil
}

func getMaxArg(args []string, logger zerolog.Logger, tx telebot.Context) int {
	var count int

	var err error

	if len(args) > 0 {
		count, err = strconv.Atoi(args[0])

		if err != nil {
			msg := fmt.Sprintf("%s não é um número válido", args[0])
			logger.Warn().Msg(msg)
			_ = tx.Reply(msg)
			count = defaultCountResults
		}

		if count > maxResults {
			msg := fmt.Sprintf("O número máximo de resultados é %v", maxResults)
			logger.Warn().Msg(msg)

			_ = tx.Reply(msg)
			count = maxResults
		}
	} else {
		count = defaultCountResults
	}

	return count
}
