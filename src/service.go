package src

import (
	"time"
	"github.com/fragmenta/server"
	"github.com/techfront/core/src/kernel/schedule"
	"github.com/techfront/core/src/component/topic/service"
	"github.com/techfront/core/src/lib/twitter"
	"github.com/techfront/core/src/lib/vk"
)

/**
 * Инициализация и конфигурирование сервисов.
 */
func setupService(server *server.Server) {
	if !server.Production() {
		return
	}

	config := server.Configuration()
	context := schedule.NewContext()
	now := time.Now().UTC()

	if config["sparkpost_key"] != "" {
		digestInterval := 55 * time.Hour
		digestTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 0, 0, time.UTC)

		schedule.At(topicservice.SendDigest, context, digestTime, digestInterval)
	}

	if config["vk_access_token"] != "" {
		vk.Setup(config["vk_access_token"], config["vk_group_id"])

		// vkTime := time.Date(now.Year(), now.Month(), now.Day(), 12, 55, 0, 0, time.UTC)
		vkInterval := 2 * time.Hour

		/**
		 * Публикация поста пре деплое:
		 * vkTime := now.Add(time.Second * 5)
		 */
		vkTime := now.Add(time.Second * 5)

		schedule.At(topicservice.VkPostTopTopic, context, vkTime, vkInterval)
	}

	if config["twitter_secret"] != "" {
		twitter.Setup(config["twitter_key"], config["twitter_secret"], config["twitter_token"], config["twitter_token_secret"])

		// tweetTime := time.Date(now.Year(), now.Month(), now.Day(), 11, 0, 0, 0, time.UTC)
		tweetInterval := 2 * time.Hour

		/**
		 * Публикация поста пре деплое:
		 * tweetTime := now.Add(time.Second * 5)
		 */
		tweetTime := now.Add(time.Second * 5)

		schedule.At(topicservice.TweetTopTopic, context, tweetTime, tweetInterval)
	}
}
