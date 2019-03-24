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
		Name:    "Help",
		Desc:    "List available commands",
		Example: "/help@@",
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

	if akita.SimpleCommandArgs("help", message.Text)[0] == "botfather" {
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
		b.WriteString("<b>Loaded commands:\n</b>")
		l := b.Len()

		for _, command := range cmd.Akita.Commands {
			help := command.Command.Help()

			if help.Desc == "" {
				continue
			}

			b.WriteString("\n")
			b.WriteString(help.HTMLString(true))
		}

		if b.Len() == l {
			b.WriteString("<i>no commands with descriptions loaded</i>\n")
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.String())
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "HTML"
	return cmd.Akita.SendMessage(msg)
}
