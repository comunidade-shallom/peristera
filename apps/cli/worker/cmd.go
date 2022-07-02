package worker

import (
	"fmt"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var Worker = &cli.Command{
	Name:  "worker",
	Usage: "Start telegram bot worker",
	Action: func(c *cli.Context) error {
		cfg := config.Ctx(c.Context)
		logger := zerolog.Ctx(c.Context).
			With().
			Str("context", "worker").
			Logger()

		logger.Debug().Msg("Creating bot instance")

		bot, err := telegram.NewBot(c.Context, *cfg)

		if err != nil {
			logger.Error().Err(err).Msg("Starting error")

			return err
		}

		ctx, _ := support.WithKillSignal(c.Context)

		go func() {
			<-ctx.Done()
			logger.Warn().Msg("Stoping bot...")
			bot.Stop()
			logger.Warn().Msg("Stoped...")
		}()

		logger.Info().Msg("Starting bot...")

		bot.OnError = func(err error, tctx telebot.Context) {
			_ = tctx.Reply(fmt.Sprintf("Error: %s", err.Error()))

			logger.Error().Err(err).Msg("Bot error")
		}

		bot.Start()

		logger.Info().Msg("Done")

		return nil
	},
}
