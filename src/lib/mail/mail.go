package mail

import (
	"github.com/techfront/core/src/kernel/view"
	sparkpost "github.com/SparkPost/gosparkpost"
)

var from sparkpost.From
var sparkpostClient sparkpost.Client

func Setup(config map[string]string) error {
	from = sparkpost.From{
		Email: config["mail_from"],
		Name: config["project_name"],
	}

	sparkpostConfig := &sparkpost.Config{
		BaseUrl: "https://api.sparkpost.com",
		ApiKey: config["sparkpost_key"],
		ApiVersion: 1,
	}
	err := sparkpostClient.Init(sparkpostConfig)

	return err
}

func Send(recipients []string, subject string, template string, variables map[string]interface{}, substitutionData map[string]map[string]interface{}) error {
	v := view.New(template)
	v.Templates = []string{template}
	v.Vars = variables
	html, err := v.RenderToString()
	if err != nil {
		return err
	}

	var recipientList []sparkpost.Recipient
	for _, v := range recipients {
		recipientList = append(recipientList, sparkpost.Recipient{
			Address: sparkpost.Address{Email: v},
		})
	}

	tx := &sparkpost.Transmission{
		Recipients: recipientList,
		SubstitutionData: substitutionData,
		Content: sparkpost.Content{
			HTML: html,
			From: from,
			Subject: subject,
		},
	}

	_, _, err = sparkpostClient.Send(tx)
	if err != nil {
		return err
	}

	return nil
}