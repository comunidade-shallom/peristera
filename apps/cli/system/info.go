package system

import (
	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support/system"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var InfoCmd = &cli.Command{
	Name:  "info",
	Usage: "load system info",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "notify",
			Usage:       "send information to admins",
			DefaultText: "false",
		},
	},
	Action: func(cmd *cli.Context) error {
		cfg := *config.Ctx(cmd.Context)

		data, err := system.New()
		if err != nil {
			return err
		}

		if cmd.Bool("notify") {
			admins := cfg.Telegram.Admins

			logger := zerolog.Ctx(cmd.Context).
				With().
				Str("context", "system").
				Logger()

			if len(admins) == 0 {
				return NoAdminsDefined
			}

			bot, err := telegram.NewBot(cfg)
			if err != nil {
				return err
			}

			msg := data.MarkdownV2(cmd.Args().First())

			for _, id := range admins {
				_, err = bot.Send(&telebot.User{
					ID: id,
				}, msg, telebot.ModeMarkdownV2)

				if err != nil {
					return err
				}
			}

			logger.Info().Msgf("System info sent do telegram %v", admins)

			return nil
		}

		return data.Println()
	},
}
