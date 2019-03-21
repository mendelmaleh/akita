package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mendelmaleh/akita"
	"strings"
)

func init() {
	akita.RegisterCommand(&helpCommand{})
}

type helpCommand struct {
	akita.CommandBase
}

func (cmd helpCommand) Help() akita.Help {
	return akita.Help{
		Name: "Help",
		Botfather: [][]string{
			[]string{"help", "get a list of available commands"},
		},
	}
}

func (cmd helpCommand) Trigger(message tgbotapi.Message) bool {
	return akita.SimpleCommand("help", message.Text)
}

func (cmd helpCommand) Exec(message tgbotapi.Message) error {
	var b strings.Builder

	if message.CommandArguments() == "botfather" {
		for k, command := range cmd.Akita.Commands {
			help := command.Command.Help().BotfatherString()

			if help != "" {
				b.WriteString(help)
				if k+1 != len(cmd.Akita.Commands) {
					b.WriteString("\n")
				}
			}
		}
	} else {
		b.WriteString("*Loaded commands:\n*")
		l := b.Len()

		for _, command := range cmd.Akita.Commands {
			help := command.Command.Help()

			if help.Desc == "" {
				continue
			}

			b.WriteString("\n")
			b.WriteString(help.String(true))
		}

		if b.Len() == l {
			b.WriteString("_no commands with descriptions_\n")
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.String())
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "Markdown"
	return cmd.Akita.SendMessage(msg)
}
