package topicservice

import (
	"fmt"
	"time"
	"strings"
	"github.com/fragmenta/query"
	"github.com/techfront/core/src/kernel/schedule"
	"github.com/techfront/core/src/lib/mail"
	"github.com/techfront/core/src/lib/mailchimp"
	"github.com/techfront/core/src/component/topic"
)

func SendDigest(context schedule.Context) {
	context.Log("#info Sending digest")

	q := topic.Popular()
	q.Where("topic_status >= 100")
	q.Where("topic_created_at > current_timestamp - interval '14 day'")
	q.Where("topic_newsletter_at IS NULL")
	q.Order("rank(score(topic_count_upvote, topic_count_downvote, topic_count_flag), topic_created_at) desc, topic_id desc")
	q.Limit(10)

	topicList, err := topic.FindAll(q)
	if err != nil {
		context.Logf("#error %s", err)
		return
	}

	if len(topicList) < 7 {
		context.Logf("#warn no topics found for newsletter")
		return
	}

	recipients, err := mailchimp.FindAllSubscribers()
	if err != nil {
		context.Logf("#error %s", err)
		return
	}

	// For testing
	// recipients = []string{"alxshelepenok@gmail.com"}

	offer := context.Config("digest_offer_code")
	subject := fmt.Sprintf("Дайджест: \"%s\" и многое другое", strings.TrimSpace(topicList[0].Name))
	variables := map[string]interface{}{"topics": topicList, "digest_offer_code": offer}
	substitutionData := map[string]map[string]interface{}{"dynamic_html": {
		"link_unsubscribe": "<a href=\"https://techfront.org/unsubscribe?email={{{address.email}}}\" target=\"_blank\" style=\"color: #7f8c8d;text-decoration: underline;font-family: -apple-system,system-ui,BlinkMacSystemFont,&quot;Segoe UI&quot;,Roboto,&quot;Helvetica Neue&quot;,Arial,sans-serif;-webkit-box-sizing: border-box;box-sizing: border-box;\">отписаться от рассылки</a>",
	}}
	err = mail.Send(recipients, subject, "component/user/template/mail/digest", variables, substitutionData)
	if err != nil {
		context.Logf("#error %s", err)
		return
	}

	params := map[string]string{"topic_newsletter_at": query.TimeString(time.Now().UTC())}
	for _, topicEntity := range topicList {
		err = topicEntity.Update(params)
		if err != nil {
			context.Logf("#error %s", err)
			return
		}
	}
}
