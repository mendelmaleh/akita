package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mendelmaleh/akita"
)

func init() {
	akita.RegisterCommand(&startCommand{})
}

type startCommand struct {
	akita.CommandBase
}

func (cmd startCommand) Help() akita.Help {
	return akita.Help{
		Name: "Start",
		Botfather: [][]string{
			[]string{"start", "start the bot"},
		},
	}
}

func (cmd startCommand) Trigger(message tgbotapi.Message) bool {
	return akita.SimpleCommand("start", message.Text)
}

func (cmd startCommand) Exec(message tgbotapi.Message) error {
	val := cmd.Config.Get("start")
	var text string
	if val == nil {
		text = "sup?"
	} else {
		text = val.(string)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	return cmd.Akita.SendMessage(msg)
}
