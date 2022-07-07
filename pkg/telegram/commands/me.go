package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/comunidade-shallom/peristera/pkg/support"
	"gopkg.in/telebot.v3"
)

func (h Commands) Me(tx telebot.Context) error {
	var builder strings.Builder

	addUser := func(user *telebot.User) {
		if user == nil {
			return
		}

		builder.WriteString("\n*Name: *" + support.AddSlashes(user.FirstName+" "+user.LastName))
		builder.WriteString("\n*Username: *`" + support.AddSlashes(user.Username) + "`")
		builder.WriteString("\n*User ID: * `" + strconv.Itoa(int(user.ID)) + "`")
	}

	addChat := func(chat *telebot.Chat) {
		builder.WriteString("\n*Chat ID: * `" + strconv.Itoa(int(chat.ID)) + "`")
		builder.WriteString("\n*Chat Title: * " + support.AddSlashes(chat.Title))
		builder.WriteString("\n*Chat Name: * " + support.AddSlashes(chat.FirstName+" "+chat.LastName))
		builder.WriteString("\n*Chat Type: * " + string(chat.Type))
		builder.WriteString(fmt.Sprintf("\n*Chat Private: * %v", chat.Private))
	}

	addData := func(msg *telebot.Message) {
		addUser(msg.Sender)

		if msg.Chat != nil {
			builder.WriteString("\n\n*\\- Chat Data:*\n")
			addChat(msg.Chat)
		}

		if msg.OriginalSender != nil {
			builder.WriteString("\n\n*\\- Original Sender:*\n")
			addUser(msg.OriginalSender)
		}

		if msg.OriginalChat != nil {
			builder.WriteString("\n\n*\\- Original Chat:*\n")
			addChat(msg.OriginalChat)
		}
	}

	original := tx.Message()

	builder.WriteString("*0️⃣ Message Data:*")

	addData(original)

	if original.ReplyTo != nil {
		builder.WriteString("\n\n* 1️⃣ Reply data:*")
		addData(original.ReplyTo)
	}

	return tx.Send(builder.String(), telebot.ModeMarkdownV2)
}
