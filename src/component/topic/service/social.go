package topicservice

import (
	"fmt"
	"time"

	"github.com/fragmenta/query"
	"github.com/kennygrant/sanitize"
	"github.com/techfront/core/src/kernel/schedule"
	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/lib/facebook"
	"github.com/techfront/core/src/lib/twitter"
	"github.com/techfront/core/src/lib/urlshortener"
	"github.com/techfront/core/src/lib/vk"
)

/**
* TweetTopTopic tweets the top topic
 */
func TweetTopTopic(context schedule.Context) {
	context.Log("#info Sending top topic tweet")

	// Get the top topic which has not been tweeted yet, newer than 1 day (we don't look at older topics)
	q := topic.Popular().Limit(1).Order("rank(score(topic_count_upvote, topic_count_downvote, topic_count_flag), topic_created_at) desc, score(topic_count_upvote, topic_count_downvote, topic_count_flag) desc, topic_id desc")

	// Fetch published only
	q.Where("topic_status >= 100")

	// Don't fetch old topics - at some point soon this can come down to 1 day
	// as all older topics will have been tweeted
	q.Where("topic_created_at > current_timestamp - interval '12 hours'")

	// Don't fetch topics that have already been tweeted
	q.Where("topic_tw_posted_at IS NULL")

	// Получение топиков
	results, err := topic.FindAll(q)
	if err != nil {
		context.Logf("#error getting top topic tweet %s", err)
		return
	}

	if len(results) > 0 {
		topicEntity := results[0]
		url := topicEntity.Url

		if topicEntity.IsVideo() || topicEntity.IsQuestion() || url == ""  {
			return
		}

		if len(topicEntity.Text) > 0 {
			url =  fmt.Sprintf("https://techfront.org/topics/%d", topicEntity.Id)
		}

		tweet := fmt.Sprintf("%s %s", topicEntity.Name, url)
		_, err := twitter.Tweet(tweet)
		if err != nil {
			context.Logf("#error tweeting top topic %s", err)
			return
		}

		// Record that this topic has been tweeted in db
		params := map[string]string{"topic_tw_posted_at": query.TimeString(time.Now().UTC())}
		err = topicEntity.Update(params)
		if err != nil {
			context.Logf("#error updating top topic tweet %s", err)
			return
		}
	} else {
		context.Logf("#warn no top topic found for tweet")
	}

}

func VkPostTopTopic(context schedule.Context) {
	context.Log("#info posting top topic vk")

	// Get the top topic
	q := topic.Popular().Limit(1).Order("rank(score(topic_count_upvote, topic_count_downvote, topic_count_flag), topic_created_at) desc, score(topic_count_upvote, topic_count_downvote, topic_count_flag) desc, topic_id desc")

	// Fetch published only
	q.Where("topic_status >= 100")

	// Don't fetch old topics
	q.Where("topic_created_at > current_timestamp - interval '12 hours'")

	q.Where("topic_vk_posted_at IS NULL")

	// Получение топиков
	results, err := topic.FindAll(q)
	if err != nil {
		context.Logf("#error getting top topic for vk %s", err)
		return
	}

	if len(results) > 0 {
		topicEntity := results[0]
		context.Logf("#info vk posting %s", topicEntity.Name)

		longUrlTopic := fmt.Sprintf("https://techfront.org/topics/%d", topicEntity.Id)
		shortUrlTopic, err := urlshortener.Shorten(longUrlTopic)
		if err != nil {
			context.Logf("#error while shorting url %s", err)
			return
		}

		if topicEntity.IsVideo() || topicEntity.IsQuestion() || topicEntity.Url == ""  {
			return
		}

		thumbnail := ""
		message := fmt.Sprintf("%s: %s\n\nКомментарии: %s", topicEntity.Name, topicEntity.Url, shortUrlTopic)
		urlSource := topicEntity.Url
		text := sanitize.HTML(topicEntity.Text)

		if len(text) > 0 {
			if len(topicEntity.Thumbnail) > 0 {
				thumbnail = fmt.Sprintf("https://techfront.org/uploads%s", topicEntity.Thumbnail)
			}

			message = fmt.Sprintf("%s: %s\n\n%s", topicEntity.Name, shortUrlTopic, text)
			urlSource = ""
		}

		if err := vk.Post(message, thumbnail, urlSource); err != nil {
			context.Logf("#error vk post top topic %s", err)
			return
		}

		params := map[string]string{"topic_vk_posted_at": query.TimeString(time.Now().UTC())}
		if err = topicEntity.Update(params); err != nil {
			context.Logf("#error updating top topic vk %s", err)
			return
		}
	} else {
		context.Logf("#warn no top topic found for vk")
	}
}

/**
* FacebookPostTopTopic facebook posts the top topic
 */
func FacebookPostTopTopic(context schedule.Context) {
	context.Log("#info posting top topic facebook")

	// Построение запроса
	q := topic.Popular().Limit(1).Order("rank(score(topic_count_upvote, topic_count_downvote, topic_count_flag), topic_created_at) desc, score(topic_count_upvote, topic_count_downvote, topic_count_flag) desc, topic_id desc")

	q.Where("topic_status >= 100")

	q.Where("topic_created_at > current_timestamp - interval '6 hours'")

	q.Where("topic_fb_posted_at IS NULL")

	// Получение топиков
	results, err := topic.FindAll(q)
	if err != nil {
		context.Logf("#error getting top topic for fb %s", err)
		return
	}

	if len(results) > 0 {
		topicEntity := results[0]
		context.Logf("#info facebook posting %s", topicEntity.Name)
		err := facebook.Post(topicEntity.Name, topicEntity.Url)
		if err != nil {
			context.Logf("#error facebook post top topic %s", err)
			return
		}

		params := map[string]string{"topic_fb_posted_at": query.TimeString(time.Now().UTC())}
		err = topicEntity.Update(params)
		if err != nil {
			context.Logf("#error updating top topic fb %s", err)
			return
		}
	} else {
		context.Logf("#warn no top topic found for fb")
	}
}
