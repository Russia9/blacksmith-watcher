package commands

import (
	"blacksmith-watcher/types"
	"blacksmith-watcher/utils"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

func SendOfferNotification(userID int, shop types.BlacksmithShop, offer types.BlacksmithOffer) {

}

func SendShopNotification(userID int, shop types.BlacksmithShop) {
	chat, err := utils.Bot.ChatByID(strconv.Itoa(userID))
	if err == nil {
		var notification strings.Builder
		notification.WriteString("<b>Открылась лавка, на которую вы подписаны:</b>\n")
		notification.WriteString(shop.Kind + " ")
		notification.WriteString("<a href=\"http://t.me/share/url?url=/ws_" + shop.Link + "\">«" + shop.Name + "»</a>")
		notification.WriteString(" игрока ")
		notification.WriteString(shop.OwnerCastle)
		if shop.OwnerTag != nil {
			notification.WriteString("[" + *shop.OwnerTag + "]")
		}
		notification.WriteString(shop.OwnerName)

		GoToShopButton := utils.UnsubscribeShop.URL("Перейти в лавку", "http://t.me/share/url?url=/ws_"+shop.Link)
		utils.UnsubscribeShop.Inline(utils.UnsubscribeShop.Row(GoToShopButton), utils.UnsubscribeShop.Row(utils.UnsubscribeShopButton))
		_, err = utils.Bot.Send(chat, notification.String(), utils.UnsubscribeShop, telebot.ModeHTML)
	}
}
