package telegram

import (
	"os"

	"github.com/comunidade-shallom/peristera/pkg/support"
	"gopkg.in/telebot.v3"
)

func Document(file *os.File, caption string) *telebot.Document {
	document := &telebot.Document{File: telebot.FromDisk(file.Name())}
	document.Caption = support.AddSlashes(caption)
	document.FileName = file.Name()

	return document
}
