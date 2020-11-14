package main

import (
	"blacksmith-watcher/commands"
	"blacksmith-watcher/types"
	"blacksmith-watcher/utils"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

func InitBot(telegramToken string, consumer *kafka.Consumer) error {
	var err error
	utils.Bot, err = telebot.NewBot(
		telebot.Settings{
			Token:  telegramToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
	if err != nil {
		return err
	}

	utils.Bot.Handle("/start", commands.StartCommandHandler)

	utils.Bot.Handle("/subscribes", commands.SubscribesCommandHandler)
	utils.Bot.Handle("/subs", commands.SubscribesCommandHandler)

	utils.Bot.Handle("/subscribe", commands.SubscribeCommandHandler)
	utils.Bot.Handle("/sub", commands.SubscribeCommandHandler)
	utils.Bot.Handle(&utils.SubscribeButton, commands.UnsubscribeCommandHandler)

	utils.Bot.Handle("/unsubscribe", commands.UnsubscribeCommandHandler)
	utils.Bot.Handle("/unsub", commands.UnsubscribeCommandHandler)
	utils.Bot.Handle(&utils.UnsubscribeButton, commands.UnsubscribeCommandHandler)

	utils.Bot.Handle(telebot.OnText, commands.ItemHandler)

	go utils.Bot.Start()

	err = consumer.SubscribeTopics([]string{"cw3-yellow_pages"}, nil)
	if err != nil {
		return err
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var message []types.BlacksmithShop
			err = json.Unmarshal([]byte(msg.Value), &message)
			if err != nil {
				sentry.CaptureException(err)
				log.Err(err).Send()
			}

			for _, shop := range message { // Going through all shops
				if time.Now().Unix()-310 > shop.LastOpenTime.Unix() {

					//for _, offer := range shop.Offers {
					//
					//}
				}
				shop.LastOpenTime = time.Now()
				types.UpdateShop(shop)
			}

			if err != nil {
				sentry.CaptureException(err)
				log.Error().Err(err).Send()
			}
		} else {
			sentry.CaptureException(err)
			log.Error().Err(err).Send()
		}
	}
}
