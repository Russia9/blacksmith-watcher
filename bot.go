package main

import (
	"blacksmith-watcher/commands"
	"blacksmith-watcher/types"
	"blacksmith-watcher/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
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
			err = json.Unmarshal(msg.Value, &message)
			if err != nil {
				sentry.CaptureException(err)
				log.Err(err).Send()
			}
			fmt.Println(1)

			for _, shop := range message {
				if time.Now().Unix()-310 > shop.LastOpenTime.Unix() {
					var users = make([]types.User, 0)
					dbResult, err := utils.DB.Collection("users").Find(context.TODO(), bson.M{"shopsubs": bson.M{"$in": []string{shop.Link}}})
					if err != nil {
						log.Error().Err(err).Send()
						continue
					}
					dbResult.Decode(&users)
					fmt.Println(2, len(users))
					if err != nil {
						log.Error().Err(err).Send()
						continue
					}

					for _, user := range users {
						fmt.Println(2)
						if err != nil {
							log.Error().Err(err).Send()
							continue
						}

						chat, err := utils.Bot.ChatByID(strconv.Itoa(user.ID))
						if err != nil {
							log.Error().Err(err).Send()
							continue
						}
						log.Info().Str(strconv.Itoa(user.ID), shop.Link)
						utils.Bot.Send(chat, shop.Link)
					}

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
