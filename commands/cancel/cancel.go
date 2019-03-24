package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mendelmaleh/akita"
)

func init() {
	akita.RegisterCommand(&cancelCommand{})
}

type cancelCommand struct {
	akita.CommandBase
}

func (cmd cancelCommand) Help() akita.Help {
	return akita.Help{
		Name: "Cancel",
		Botfather: [][]string{
			[]string{"cancel", "cancel current command"},
		},
	}
}

func (cmd cancelCommand) Trigger(message tgbotapi.Message) bool {
	return akita.SimpleCommand("cancel", message.Text)
}

func (cmd cancelCommand) Exec(message tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "<i>current commands canceled</i>")
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      true,
	}

	for _, command := range cmd.Akita.Commands {
		command.StopWaiting(message.From.ID)
	}

	return cmd.Akita.SendMessage(msg)
}
