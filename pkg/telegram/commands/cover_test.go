//nolint:funlen
package commands_test

import (
	"testing"

	"github.com/comunidade-shallom/diakonos/pkg/covers"
	"github.com/comunidade-shallom/peristera/pkg/telegram/commands"
	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
)

func TestBuildCoverParams(t *testing.T) {
	type test struct {
		msg  *telebot.Message
		size covers.Size
		text string
		name string
	}

	tests := []test{
		{
			name: "empty payload",
			size: covers.Size{Width: 1080, Height: 1080},
			text: "",
			msg: &telebot.Message{
				Payload: "",
			},
		},
		{
			name: "single payload",
			size: covers.Size{Width: 1080, Height: 1080},
			text: "foo bar",
			msg: &telebot.Message{
				Payload: "foo bar",
			},
		},
		{
			name: "sized payload",
			size: covers.Size{Width: 200, Height: 200},
			text: "foo bar",
			msg: &telebot.Message{
				Payload: "200x200 foo bar ",
			},
		},
		{
			name: "sized payload (multiple)",
			size: covers.Size{Width: 100, Height: 200},
			text: "foo bar 50x50",
			msg: &telebot.Message{
				Payload: "100x200 foo bar 50x50",
			},
		},
		{
			name: "reply photo",
			size: covers.Size{Width: 150, Height: 200},
			text: "Bulma",
			msg: &telebot.Message{
				ReplyTo: &telebot.Message{
					Caption: "Bulma",
					Photo: &telebot.Photo{
						Width:  150,
						Height: 200,
					},
				},
			},
		},
		{
			name: "reply text",
			size: covers.Size{Width: 1080, Height: 1080},
			text: "Freeza",
			msg: &telebot.Message{
				ReplyTo: &telebot.Message{
					Text: "Freeza",
				},
			},
		},
		{
			name: "reply text (sized)",
			size: covers.Size{Width: 400, Height: 90},
			text: "Freeza",
			msg: &telebot.Message{
				Payload: "400x90",
				ReplyTo: &telebot.Message{
					Text: "Freeza",
				},
			},
		},
	}

	for _, val := range tests {
		current := val

		t.Run(current.name, func(t *testing.T) {
			size, text := commands.BuildCoverParams(current.msg)

			assert.Equal(t, current.text, text)
			assert.Equal(t, current.size, size)
		})
	}
}
