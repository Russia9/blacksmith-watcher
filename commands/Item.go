package commands

import (
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

func ItemHandler(message *telebot.Message) {
	if message.IsForwarded() {
		switch message.OriginalSender.ID {
		case 768128010:
			if strings.Contains(message.Text, "Recipe of ") {
				//item := strings.Split(strings.Split(message.Text, ":")[0], "Recipe of ")[1]
				// TODO: Search item's id
			}
		}
	} else {

	}
}
