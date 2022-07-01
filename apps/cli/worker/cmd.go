package worker

import (
	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var Worker = &cli.Command{
	Name:  "worker",
	Usage: "Start telegram bot worker",
	Action: func(ctxCli *cli.Context) error {
		cfg := config.Ctx(ctxCli.Context)

		pterm.Debug.Println("Creating bot instance")

		bot, err := telegram.NewBot(cfg)
		if err != nil {
			return err
		}

		ctx, _ := support.WithKillSignal(ctxCli.Context)

		go func() {
			<-ctx.Done()
			pterm.Warning.Println("Stoping bot...")
			bot.Stop()
			pterm.Debug.Println("Stoped...")
		}()

		pterm.Info.Println("Starting bot...")

		bot.OnError = func(err error, _ telebot.Context) {
			pterm.Error.Println(err.Error())
		}

		bot.Start()

		pterm.Info.Println("Done")

		return nil
	},
}
