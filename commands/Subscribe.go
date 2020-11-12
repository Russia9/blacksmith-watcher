package commands

import (
	"blacksmith-watcher/types"
	"blacksmith-watcher/utils"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

func SubscribeHandler(message *telebot.Message) {
	if message.IsForwarded() {
		switch message.OriginalSender.ID {
		case 768128010:
			if strings.Contains(message.Text, "Recipe of ") {
				item := strings.Split(strings.Split(message.Text, ":")[0], "Recipe of ")[1]
				
				utils.UnsubscribeItemButton.Data = item
				utils.UnsubscribeItem.Inline(utils.UnsubscribeItem.Row(utils.UnsubscribeItemButton))
				
				if types.GetUser(message.Sender.ID).SubItem(item) {
					utils.Bot.Send(message.Sender, "Подписка на предмет *"+item+"* оформлена\\.", utils.UnsubscribeItem, telebot.ModeMarkdownV2)
				} else {
					utils.Bot.Send(message.Sender, "Вы уже подписаны на предмет *"+item+"*\\.", utils.UnsubscribeItem, telebot.ModeMarkdownV2)
				}
			}
		}
	}
}