package types

import (
	"blacksmith-watcher/utils"
	"context"
	"strings"
)

type User struct {
	ID       int      `json:"id"`
	ShopSubs []string `json:"shop_subs"`
	ItemSubs []string `json:"item_subs"`
}

func (user User) SubItem(item string) bool {
	if !utils.Contains(user.ItemSubs, item) {
		user.ItemSubs = append(user.ItemSubs, strings.TrimSpace(item))
		filter := struct {
			ID int `json:"id"`
		}{ID: user.ID}
		utils.DB.Collection("users").FindOneAndReplace(context.TODO(), filter, user)
		return true
	}
	return false
}

func (user User) UnsubItem(item string) bool {
	if utils.Contains(user.ItemSubs, item) {
		user.ItemSubs = utils.Remove(user.ItemSubs, strings.TrimSpace(item))
		filter := struct {
			ID int `json:"id"`
		}{ID: user.ID}
		utils.DB.Collection("users").FindOneAndReplace(context.TODO(), filter, user)
		return true
	}
	return false
}

func GetUser(id int) User {
	var result User
	filter := struct {
		ID int `json:"id"`
	}{ID: id}
	findUserResult := utils.DB.Collection("users").FindOne(context.TODO(), filter)
	if findUserResult.Err() == nil {
		findUserResult.Decode(&result)
	} else {
		result = User{
			ID:       id,
			ShopSubs: make([]string, 0),
			ItemSubs: make([]string, 0),
		}
		utils.DB.Collection("users").InsertOne(context.TODO(), result)
	}
	return result
}
