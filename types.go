package akita

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"sync"
)

type Help struct {
	Name      string
	Desc      string
	Example   string
	Botfather [][]string
}

func (h Help) String(full bool) string {
	var b strings.Builder
	b.WriteString(h.Name)

	if full {
		b.WriteString("\n")
	} else {
		b.WriteString(" - ")
	}
	b.WriteString(h.Desc)

	if full {
		b.WriteString("\neg: ")
		b.WriteString(h.Example)
		b.WriteString("\n")
	}

	return b.String()
}

func (h Help) HTMLString(full bool) string {
	var b strings.Builder
	b.WriteString("<b>• ")
	b.WriteString(h.Name)
	b.WriteString("</b>")

	if full {
		b.WriteString("\n")
		b.WriteString("<i>")
	} else {
		b.WriteString(" - ")
	}
	b.WriteString(h.Desc)

	if full {
		b.WriteString("</i>")
		b.WriteString("\neg: ")
		b.WriteString(h.Example)
		b.WriteString("\n")
	}

	return b.String()
}

func (h Help) BotfatherString() string {
	if len(h.Botfather) == 0 {
		return ""
	}

	var b strings.Builder

	for k, v := range h.Botfather {
		b.WriteString(v[0])
		b.WriteString(" - ")
		b.WriteString(v[1])
		if k+1 != len(h.Botfather) {
			b.WriteString("\n")
		}
	}

	return b.String()
}

type Command interface {
	Help() Help
	Init(*CommandState, *Akita) error
	Trigger(tgbotapi.Message) bool
	Exec(tgbotapi.Message) error
	ExecWaiting(tgbotapi.Message) error
}

type CommandBase struct {
	*CommandState
	*Akita
}

func (CommandBase) Help() Help { return Help{} }

func (cmd *CommandBase) Init(c *CommandState, a *Akita) error {
	cmd.CommandState = c
	cmd.Akita = a

	return nil
}

func (CommandBase) Trigger(tgbotapi.Message) bool { return false }

func (CommandBase) Exec(tgbotapi.Message) error { return nil }

func (CommandBase) ExecWaiting(tgbotapi.Message) error { return nil }

func (cmd CommandBase) Get(key string) interface{} {
	return cmd.Akita.Config.Get(key)
}

func (cmd CommandBase) Set(key string, value interface{}) {
	cmd.Akita.Config[key] = value
	cmd.Akita.Config.Save()
}

type userWaitMap struct {
	mutex    *sync.Mutex
	userWait map[int]bool
}

type CommandState struct {
	Command         Command
	waitingForReply userWaitMap
}

func NewCommandState(cmd Command) *CommandState {
	return &CommandState{
		Command: cmd,
		waitingForReply: userWaitMap{
			mutex:    &sync.Mutex{},
			userWait: map[int]bool{},
		},
	}
}

func (state *CommandState) IsWaiting(user int) bool {
	state.waitingForReply.mutex.Lock()
	defer state.waitingForReply.mutex.Unlock()
	if v, ok := state.waitingForReply.userWait[user]; ok {
		return v
	}

	return false
}

func (state *CommandState) StartWaiting(user int) {
	state.waitingForReply.mutex.Lock()
	defer state.waitingForReply.mutex.Unlock()
	state.waitingForReply.userWait[user] = true
}

func (state *CommandState) StopWaiting(user int) {
	state.waitingForReply.mutex.Lock()
	defer state.waitingForReply.mutex.Unlock()
	state.waitingForReply.userWait[user] = false
}
