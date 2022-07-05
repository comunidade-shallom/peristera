package system

import (
	"bytes"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/comunidade-shallom/peristera/pkg/config"
	"github.com/comunidade-shallom/peristera/pkg/support/errors"
	"github.com/comunidade-shallom/peristera/pkg/telegram"
	"github.com/comunidade-shallom/peristera/pkg/ytube"
	"github.com/matishsiao/goInfo"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var NoAdminsDefined = errors.Business("No admins defined", "SY:001")

var specials = []rune{'.', '\'', '(', ')', '-'}

var System = &cli.Command{
	Name:  "system",
	Usage: "load system info",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "notify",
			Usage:       "send information to admins",
			DefaultText: "false",
		},
	},
	Action: func(cmd *cli.Context) error {
		cfg := config.Ctx(cmd.Context)

		if cmd.Bool("notify") {
			admins := cfg.Telegram.Admins

			logger := zerolog.Ctx(cmd.Context).
				With().
				Str("context", "system").
				Logger()

			if len(admins) == 0 {
				return NoAdminsDefined
			}

			bot, err := telegram.NewBot(cmd.Context, *cfg, ytube.Service{})
			if err != nil {
				return err
			}

			msg, err := getMessage(cmd.Args().First())
			if err != nil {
				return err
			}

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

		return renderInfo()
	},
}

func renderInfo() error {
	info, err := goInfo.GetInfo()
	if err != nil {
		return err
	}

	infos, err := pterm.DefaultTable.
		WithBoxed().
		WithData(pterm.TableData{
			{"Hostname", info.Hostname},
			{"Platform", info.Platform},
			{"CPUs", strconv.Itoa(info.CPUs)},
			{"GoOS", info.GoOS},
			{"Core", info.Core},
			{"OS", info.OS},
			{"Kernel", info.Kernel},
		}).Srender()
	if err != nil {
		return err
	}

	list, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	lines := pterm.TableData{
		{"network", "IP"},
	}

	for _, v := range list {
		lines = append(lines, []string{v.Network(), v.String()})
	}

	ips, err := pterm.DefaultTable.
		WithBoxed().
		WithHasHeader().
		WithData(lines).Srender()
	if err != nil {
		return err
	}

	panels, err := pterm.DefaultPanel.WithPanels(pterm.Panels{
		{{Data: infos}},
		{{Data: ips}},
	}).Srender()
	if err != nil {
		return err
	}

	pterm.DefaultBox.
		WithTitle("System Info").
		WithTitleBottomRight().
		WithRightPadding(0).
		WithBottomPadding(0).
		Println(panels)

	return nil
}

func getMessage(head string) (string, error) {
	list, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	info, err := goInfo.GetInfo()
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	if head != "" {
		builder.WriteString(addslashes(head))
		builder.WriteString("\n")
	}

	builder.WriteString("\n*System Details*\n")
	builder.WriteString("\n*Hostname:* " + addslashes(info.Hostname))
	builder.WriteString("\n*Platform:* " + addslashes(info.Platform))
	builder.WriteString("\n*CPUs:* " + strconv.Itoa(info.CPUs))
	builder.WriteString("\n*GoOS:* " + addslashes(info.GoOS))
	builder.WriteString("\n*Core:* " + addslashes(info.Core))
	builder.WriteString("\n*OS:* " + addslashes(info.OS))
	builder.WriteString("\n*Kernal:* " + addslashes(info.Kernel))

	builder.WriteString("\n\n*System IPs*\n")

	for _, ip := range list {
		builder.WriteString("\n`" + addslashes(ip.String()) + "`")
	}

	builder.WriteString("\n\n*System Time:* \n" + addslashes(time.Now().Format(time.RFC3339)) + "\n")

	return builder.String(), nil
}

func addslashes(str string) string {
	var buf bytes.Buffer

	for _, char := range str {
		for _, sp := range specials {
			if sp == char {
				buf.WriteRune('\\')

				continue
			}
		}

		buf.WriteRune(char)
	}

	return buf.String()
}
