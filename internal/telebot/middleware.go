package telebot

import (
	tele "gopkg.in/telebot.v3"
)

// HandleFuncWithMessage is a function type that handles a message and returns a Message and an error.
type HandleFuncWithMessage func(tele.Context) (*tele.Message, error)

// wrap is a function that wraps a HandleFuncWithMessage and adds the message to a storage.
func (b *Bot) wrap(handler HandleFuncWithMessage) tele.HandlerFunc {
	return func(c tele.Context) error {
		msg, err := handler(c)

		b.lastSentMessage.set(c.Sender().ID, msg.Text)

		return err
	}
}
