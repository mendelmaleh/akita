package akita

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
	"strings"
)

var commands []*CommandState

func RegisterCommand(cmd Command) {
	commands = append(commands, NewCommandState(cmd))
}

func SimpleCommand(command, text string) bool {
	return regexp.MustCompile("^/(" + command + ")(@\\w+)?( .+)?$").MatchString(text)
}

func SimpleArgCommand(command, text string, args int) bool {
	matches := regexp.MustCompile("^/(" + command + ")(@\\w+)?( .+)?$").FindStringSubmatch(text)
	if len(matches) < 4 {
		return false
	}

	matchedArgs := len(strings.Split(strings.Trim(matches[3], " "), " "))
	return args == matchedArgs
}

func (a *Akita) commandRouter(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if update.Message == nil {
		return
	}

	for _, cmd := range a.Commands {
		if cmd.IsWaiting(update.Message.From.ID) {
			if err := cmd.Command.ExecWaiting(*update.Message); err != nil {
				a.commandError(cmd.Command.Help().Name, *update.Message, err)
			}
			continue
		}

		if cmd.Command.Trigger(*update.Message) {
			if err := cmd.Command.Exec(*update.Message); err != nil {
				a.commandError(cmd.Command.Help().Name, *update.Message, err)
			}
		}
	}
}

func (a *Akita) commandInit() {
	for _, cmd := range a.Commands {
		err := cmd.Command.Init(cmd, a)
		if err != nil {
			log.Printf("Error starting command %s: %s\n", cmd.Command.Help().Name, err.Error())
		} else {
			log.Printf("Started command %s!", cmd.Command.Help().Name)
		}
	}
}

func (a *Akita) commandError(name string, message tgbotapi.Message, err error) {
	var msg tgbotapi.MessageConfig

	if a.API.Debug {
		msg = tgbotapi.NewMessage(message.Chat.ID, err.Error())
	} else {
		msg = tgbotapi.NewMessage(message.Chat.ID, "An error occured!")
		log.Println("Error processing command: " + err.Error())
	}

	msg.ReplyToMessageID = message.MessageID

	_, err = a.API.Send(msg)
	if err != nil {
		log.Printf("An error happened processing an error!\n%s\n", err.Error())
	}
}
