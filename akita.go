package akita

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"strings"
)

type Akita struct {
	API      *tgbotapi.BotAPI
	Config   Config
	Commands []*CommandState
}

func NewAkita(token string) *Akita {
	return NewAkitaWithClient(token, &http.Client{})
}

func NewAkitaWithClient(token string, client *http.Client) *Akita {
	bot := &Akita{}

	c, _ := LoadConfig()
	bot.Config = *c

	if token == "" {
		val := bot.Config.Get("token")
		if val == nil {
			panic("no token provided!")
		}

		token = val.(string)
	}

	api, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		panic(err)
	}

	bot.API = api
	bot.Commands = commands

	return bot
}

func (a *Akita) Start() {
	a.commandInit()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 86400

	updates, err := a.API.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	log.Println("Starting bot...")
	for update := range updates {
		go a.commandRouter(update)
	}
}

func (a *Akita) SendMessage(message tgbotapi.MessageConfig) error {
	message.Text = strings.Replace(message.Text, "@@", "@"+a.API.Self.UserName, -1)
	_, err := a.API.Send(message)
	return err
}
