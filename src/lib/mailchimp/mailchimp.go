package mailchimp

import (
	"fmt"
	"github.com/beeker1121/mailchimp-go"
	"github.com/beeker1121/mailchimp-go/lists/members"
)

/**
 * ID списка на Mailchimp, задается в конфиге.
 */
var ListID string

/**
 * Инициализация и конфигурирование.
 */
func Setup(config map[string]string) {
	ListID = config["mailchimp_list_id"]
	mailchimp.SetKey(config["mailchimp_key"])
}

/**
* Функция добавляет Email в список.
 */
func AddMember(email string) error {
	request := &members.NewParams{
		EmailAddress: email,
		Status: members.StatusSubscribed,
	}

	if _, err := members.New(ListID, request); err != nil {
		return SubscribeMember(email)
	}

	return nil
}

/**
* Функция изменяет статус пользователя на Subscribed.
*/
func SubscribeMember(email string) error {
	request := &members.UpdateParams{
		EmailAddress: email,
		Status: members.StatusSubscribed,
	}

	memberEntity, err := FindMemberByEmail(email)
	if err != nil {
		return err
	}

	if _, err := members.Update(ListID, memberEntity.ID, request); err != nil {
		return err
	}

	return nil
}

/**
* Функция изменяет статус пользователя на Unsubscribed.
*/
func UnsubscribeMember(email string) error {
	request := &members.UpdateParams{
		EmailAddress: email,
		Status: members.StatusUnsubscribed,
	}

	memberEntity, err := FindMemberByEmail(email)
	if err != nil {
		return err
	}

	if _, err := members.Update(ListID, memberEntity.ID, request); err != nil {
		return err
	}

	return nil
}

/**
* Функция возвращает объект пользователя из списка mailchimp по email.
*/
func FindMemberByEmail(email string) (*members.Member, error) {
	list, err := members.Get(ListID, nil)
	if err != nil {
		return nil, err
	}

	for _, v := range list.Members {
		if v.EmailAddress == email {
			memberEntity, err := members.GetMember(ListID, v.ID, nil)
			if err != nil {
				return nil, err
			}

			return memberEntity, nil
		} 
	}

	return nil, fmt.Errorf("#error Member \"%s\" not found in a list.", email)
}

/**
* Функция возвращает список подписанных пользователей.
*
* @return []string массив email адресов.
*/
func FindAllSubscribers() ([]string, error) {
	list, err := members.Get(ListID, nil)
	if err != nil {
		return nil, err
	}

	var emailLIst []string

	for _, v := range list.Members {
		if v.Status == members.StatusSubscribed {
			emailLIst = append(emailLIst, v.EmailAddress)
		}
	}

	return emailLIst, nil
}