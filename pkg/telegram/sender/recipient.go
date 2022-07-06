package sender

import "gopkg.in/telebot.v3"

type (
	Chats []int64
)

func (c Chats) Recipients() []telebot.Recipient {
	res := make([]telebot.Recipient, len(c))

	for index, id := range c {
		if id > 0 {
			res[index] = &telebot.User{ID: id}
		} else {
			res[index] = &telebot.Chat{ID: id}
		}
	}

	return res
}
