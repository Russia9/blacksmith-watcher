package commands

import (
	"blacksmith-watcher/types"
	"blacksmith-watcher/utils"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

func UnsubscribeHandler(callback *telebot.Callback) {
	item := strings.TrimSpace(callback.Data)

	if types.GetUser(callback.Sender.ID).UnsubItem(item) {
		utils.Bot.Send(callback.Sender, "Вы отписались от предмета *" + item + "*", telebot.ModeMarkdownV2)
	} else {
		utils.Bot.Send(callback.Sender, "Вы не подписаны на предмет *" + item + "*", telebot.ModeMarkdownV2)
	}
	utils.Bot.Respond(callback, &telebot.CallbackResponse{})
}
