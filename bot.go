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

	utils.Bot.Handle("/start", commands.StartHandler)
	utils.Bot.Handle("/subs", commands.SubsHandler)
	utils.Bot.Handle(telebot.OnText, commands.SubscribeHandler)
	utils.Bot.Handle(&utils.UnsubscribeItemButton, commands.UnsubscribeHandler)

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
				if time.Now().Unix() - 1000 > shop.LastOpenTime.Unix() {
					//commands.SendShopNotification(784726544, shop)
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
