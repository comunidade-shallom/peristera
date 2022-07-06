package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/comunidade-shallom/peristera/pkg/support"
	"gopkg.in/telebot.v3"
)

const telegramMsgMaxLen = 4000

func (h Commands) Exec(tx telebot.Context) error {
	payload := strings.TrimSpace(tx.Message().Payload)

	if strings.HasPrefix(payload, "sudo") {
		return tx.Reply("no sudo allowed")
	}

	msg, err := tx.Bot().Reply(tx.Message(), "Running...")
	if err != nil {
		return err
	}

	err = tx.Notify(telebot.Typing)
	if err != nil {
		return err
	}

	logger := h.logger(tx)

	logger.Info().Msg("Executing command...")

	args := strings.Fields(payload)
	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec

	out, err := cmd.CombinedOutput()

	result := string(out)

	if len(out) > telegramMsgMaxLen {
		logger.Warn().Msgf("Output ir more than %v", telegramMsgMaxLen)

		result = support.TruncateString(result, telegramMsgMaxLen)
		result += "\n\n--truncated--"
	}

	if err != nil {
		logger.Error().Err(err).Msg("Command failed")

		_, err = tx.Bot().Edit(
			msg,
			fmt.Sprintf("*Error:* __%s__ \n *Output:* \n ```%s```", err.Error(), result),
			telebot.ModeMarkdownV2,
		)

		return err
	}

	logger.Info().Msg("Command done")

	_, err = tx.Bot().Edit(
		msg,
		fmt.Sprintf("*Result:* \n ```%s```", result),
		telebot.ModeMarkdownV2,
	)

	return err
}
