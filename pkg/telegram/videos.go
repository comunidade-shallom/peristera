package telegram

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

const loadTimeout = 5
const defaultCountResults = 2
const maxResults = 5

func (h Handler) Videos(tx telebot.Context) error { //nolint:funlen
	_ctx, cancel := context.WithTimeout(context.Background(), time.Second*loadTimeout)

	logger := tx.Get(loggerKey).(zerolog.Logger) //nolint:forcetypeassert

	defer cancel()

	args := tx.Args()

	var count int

	var err error

	if len(args) > 0 {
		count, err = strconv.Atoi(args[0])

		if err != nil {
			_ = tx.Reply(fmt.Sprintf("%s não é um número válido", args[0]))
			count = defaultCountResults
		}

		if count > maxResults {
			_ = tx.Reply(fmt.Sprintf("O número máximo de resultados é %v", maxResults))
			count = maxResults
		}
	} else {
		count = defaultCountResults
	}

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
