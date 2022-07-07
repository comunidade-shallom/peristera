package telegram

import (
	"os"

	"gopkg.in/telebot.v3"
)

func Document(file *os.File, caption string) *telebot.Document {
	document := &telebot.Document{File: telebot.FromDisk(file.Name())}
	document.Caption = caption
	document.FileName = file.Name()

	return document
}
