package user

import (
	"github.com/fragmenta/query"
	"github.com/techfront/core/src/kernel/validate"
)

type Contact struct {
	Id int64
	UserId int64
	Name string
	Value string
}

const (
	CONTACT_NAME_PREFIX = "user_contact_"
)

var (
	ContactNameEmail string = "email"
	ContactNamePhone string = "phone"
	ContactNameSite string = "website"
	ContactNameGithub string = "github"
	ContactNameTelegram string = "telegram"
	ContactNameWhatsapp string = "whatsapp"
	ContactNameViber string = "viber"
	ContactNameSkype string = "skype"
	ContactNameVkontakte string = "vkontakte"
	ContactNameTwitter string = "twitter"
	ContactNameYoutube string = "youtube"
	ContactNameFacebook string = "facebook"
)

/**
* Функция GetFieldName получает name для поля.
*
* @return string
*/
func (m Contact) GetFieldName() string {
	return CONTACT_NAME_PREFIX + m.Name
}

/**
* Функция GetLinkHref получает href для ссылки контакта.
*
*/
func (m Contact) GetLinkHref() string {
	var href string

	switch m.Name {
		case ContactNameEmail:
			href = "mailto:" + m.Value
		case ContactNamePhone:
			href = "tel:" + m.Value
		case ContactNameSite:
			href = m.Value
		case ContactNameGithub:
			href = "https://github.com/" + m.Value
		case ContactNameTelegram:
			href = "tg://resolve?domain=" + m.Value
		case ContactNameWhatsapp:
			href = "whatsapp://send/?phone=" + m.Value
		case ContactNameViber:
			href = "viber://add?number=" + m.Value
		case ContactNameSkype:
			href = "skype:" + m.Value
		case ContactNameVkontakte:
			href = "https://vk.com/" + m.Value
		case ContactNameTwitter:
			href = "https://twitter.com/" + m.Value
		case ContactNameFacebook:
			href = "https://facebook.com/" + m.Value
		case ContactNameYoutube:
			href = "https://youtube.com/" + m.Value
		default:
			href = m.Value
	}

	return href
}

/**
* Функция GetLinkTitle получает заголовок для ссылки контакта.
*
* @return string
*/
func (m Contact) GetLinkTitle() string {
	title := m.Name + ": " + m.Value

	return title
}

/**
* Функция GetIconClass получает класс иконки для контакта.
*
* @return string
*/
func (m Contact) GetIconClass() string {
	var class string

	switch m.Name {
		case ContactNameEmail:
			class = "icon-mail"
		case ContactNamePhone:
			class = "icon-phone"
		case ContactNameSite:
			class = "icon-link-1"
		case ContactNameGithub:
			class = "icon-github"
		case ContactNameTelegram:
			class = "icon-paper-plane"
		case ContactNameWhatsapp:
			class = "icon-whatsapp"
		case ContactNameViber:
			class = "icon-phone"
		case ContactNameSkype:
			class = "icon-skype"
		case ContactNameVkontakte:
			class = "icon-vkontakte"
		case ContactNameTwitter:
			class = "icon-twitter"
		case ContactNameFacebook:
			class = "icon-facebook"
		case ContactNameYoutube:
			class = "icon-youtube-play"
		default:
			class = "icon-link-1"
	}

	return class
}

/**
* Функция IsSelected проверяет выбран ли контакт.
*
* @param n string название контакта.
* @return bool
*/
func (m Contact) IsSelected(n string) bool {
	if m.Name != n {
		return false
	}

	return true
}

/**
* Функция AllowedContactParams определяет разрешенные параметры.
*
* @return []string
*/
func AllowedContactParams() []string {
	return []string{
		CONTACT_NAME_PREFIX + ContactNameEmail,
		CONTACT_NAME_PREFIX + ContactNamePhone,
		CONTACT_NAME_PREFIX + ContactNameSite,
		CONTACT_NAME_PREFIX + ContactNameGithub,
		CONTACT_NAME_PREFIX + ContactNameTelegram,
		CONTACT_NAME_PREFIX + ContactNameWhatsapp,
		CONTACT_NAME_PREFIX + ContactNameViber,
		CONTACT_NAME_PREFIX + ContactNameSkype,
		CONTACT_NAME_PREFIX + ContactNameVkontakte,
		CONTACT_NAME_PREFIX + ContactNameTwitter,
		CONTACT_NAME_PREFIX + ContactNameYoutube,
		CONTACT_NAME_PREFIX + ContactNameFacebook,
	}
}

/**
* Функция DestroyContact удаляет контакт из базы.
*
* @return error
 */
func (m Contact) DestroyContact() error {
	return query.New("tf_user_contact", "user_contact_id").Where("user_contact_id = ?", m.Id).Delete()
}


/**
* Функция AddContact добавляет контакт пользователя в базу.
*
* @param params map[string]string параметры.
* @return int64 id записи.
* @return error 
*/
func CreateContact(params map[string]string) (int64, error) {
	return query.New("tf_user_contact", "user_contact_id").Insert(params)
}

/**
* Фукция ValidateContact проверяет значение контакта.
*
* @param param string
* @return error
*/
func ValidateContact(param string) error {
	return validate.Length(param, 2, 150)
}

/**
* Функция GetContacts получает контакты пользователя.
*
* @return []Contact
*/
func (m *User) FindAllContacts() []Contact {
	q := query.New("tf_user_contact", "user_contact_id")
	q.Where("user_contact_id_user = ?", m.Id)

	results, err := q.Results()
	if err != nil {
		return nil
	}

	var contacts []Contact
	for _, cols := range results {
		contactEntity := Contact{}
		contactEntity.Id = validate.Int(cols["user_contact_id"])
		contactEntity.UserId = validate.Int(cols["user_contact_id_user"])
		contactEntity.Name = validate.String(cols["user_contact_name"])
		contactEntity.Value = validate.String(cols["user_contact_value"])

		contacts = append(contacts, contactEntity)
	}

	return contacts
}