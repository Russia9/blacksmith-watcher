package commands

import (
	"blacksmith-watcher/utils"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"gopkg.in/tucnak/telebot.v2"
)

func StartHandler(message *telebot.Message) {
	if !message.Private() {
		return
	}

	_, err := utils.Bot.Send(message.Chat, "Привет!")
	if err != nil {
		log.Warn().Err(err).Send()
		sentry.CaptureException(err)
	}
}