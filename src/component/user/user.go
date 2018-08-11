package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/techfront/core/src/kernel/model"
	"github.com/techfront/core/src/kernel/validate"
	"github.com/fragmenta/query"

	"github.com/techfront/core/src/kernel/auth"
	"github.com/techfront/core/src/kernel/router"

	"github.com/techfront/core/src/lib/limiter"
	"github.com/techfront/core/src/lib/status"
	"github.com/techfront/core/src/lib/mail"
)

type User struct {
	model.Model
	status.ModelStatus
	VisitedAt time.Time
	Email string
	Name string
	FullName string
	EncryptedPassword string
	CreateToken string
	ResetToken string
	Avatar  string
	Text string
	Gender int64
	Contacts []Contact
	Role int64
	Power float64
	Score float64
	CommentCount int64
	TopicCount int64
	OfferCount int64
}

/**
* AllowedParams returns an array of acceptable params in update.
*/
func AllowedParams() []string {
	return []string{"user_name", "user_gender", "user_fullname", "user_text", "user_email", "user_title", "user_password"}
}

/**
* AllowedParamsAdmin returns an array of acceptable params in update for admins.
*/
func AllowedParamsAdmin() []string {
	return []string{"user_name", "user_gender", "user_fullname", "user_text", "user_email", "user_avatar", "user_status", "user_role", "user_password"}
}

/**
* AllowedUpdateParams returns an array of acceptable params in update.
*/
func AllowedUpdateParams() []string {
	return []string{"user_text", "user_gender", "user_avatar", "user_fullname", "user_email", "user_password"}
}

/**
* NewWithColumns creates a user from database columns - used by query in creating objects.
*/
func NewWithColumns(cols map[string]interface{}) *User {
	u := New()
	u.Id = validate.Int(cols["user_id"])
	u.CreatedAt = validate.Time(cols["user_created_at"])
	u.UpdatedAt = validate.Time(cols["user_updated_at"])
	u.VisitedAt = validate.Time(cols["user_visited_at"])
	u.Name = validate.String(cols["user_name"])
	u.FullName = validate.String(cols["user_fullname"])
	u.Email = validate.String(cols["user_email"])
	u.EncryptedPassword = validate.String(cols["user_encrypted_password"])
	u.Avatar = validate.String(cols["user_avatar"])
	u.Text = validate.String(cols["user_text"])
	u.Status = validate.Int(cols["user_status"])
	u.Gender = validate.Int(cols["user_gender"])
	u.Role = validate.Int(cols["user_role"])
	u.Score = validate.Float(cols["user_score"])
	u.Power = validate.Float(cols["user_power"])
	u.ResetToken = validate.String(cols["user_reset_token"])
	u.CreateToken = validate.String(cols["user_create_token"])
	u.CommentCount = validate.Int(cols["user_count_comment"])
	u.TopicCount = validate.Int(cols["user_count_topic"])
	u.OfferCount = validate.Int(cols["user_count_offer"])


	// Получение контактов
	u.Contacts = u.FindAllContacts()

	return u
}

/**
* New sets up a new user with default values.
*/
func New() *User {
	u := &User{}
	u.Model.Init()
	u.TableName = "tf_user"
	u.KeyName = "user_id"
	u.Status = status.Published
	u.Text = ""
	return u
}

/**
* Create inserts a new user.
*/
func Create(params map[string]string) (int64, error) {
	err := validateParams(params)
	if err != nil {
		return 0, err
	}

	// Check that this user email is not already in use
	if len(params["user_email"]) > 0 {
		// Try to fetch a user by this email from the db - we don't allow duplicates
		count, err := Query().Where("user_email=?", params["user_email"]).Count()
		if err != nil {
			return 0, err
		}

		if count > 0 {
			return 0, errors.New("A username with this email already exists, sorry.")
		}

	}

	// Update/add some params by default
	params["user_created_at"] = query.TimeString(time.Now().UTC())
	params["user_updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

/**
* Query a new query relation referencing this model, optionally setting a default order.
*/
func Query() *query.Query {
	p := New()
	return query.New(p.TableName, p.KeyName)
}

/**
* Where is a shortcut for the common where query on users.
*/
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

/**
* Find fetches a single record by id.
*/
func Find(id int64) (*User, error) {
	result, err := Query().Where("user_id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

/**
* FindEmail fetches a single record by email.
*/
func FindEmail(email string) (*User, error) {
	result, err := Query().Where("user_email=?", email).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

/**
* First fetches the first result for this query.
*/
func First(q *query.Query) (*User, error) {
	result, err := q.FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

/**
* FindAll fetches all results for this query.
*/
func FindAll(q *query.Query) ([]*User, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of pages constructed from the results
	var userList []*User
	for _, r := range results {
		user := NewWithColumns(r)
		userList = append(userList, user)
	}

	return userList, nil
}

/**
* Exists checks whether a user email exists.
*/
func Exists(e string) bool {
	count, err := Query().Where("user_email=?", e).Count()
	if err != nil {
		return true // default to true on error
	}

	return (count > 0)
}

/**
* validateParams these parameters conform to AcceptedParams, and pass validation
*/
func validateParams(unsafeParams map[string]string) error {
	if len(unsafeParams["user_name"]) > 0 {
		if err := CheckName(unsafeParams["user_name"]); err != nil {
			return router.BadRequestError(err)
		}

		if err := validate.Length(unsafeParams["user_name"], 2, 100); err != nil {
			return router.BadRequestError(err)
		}
	}

	if len(unsafeParams["user_email"]) > 0 {
		err := validate.Length(unsafeParams["user_email"], 3, 100)
		if err != nil {
			return router.BadRequestError(err)
		}
	}

	// Password may be blank
	if len(unsafeParams["user_password"]) > 0 {
		// Report error for length between 0 and 8 chars
		err := validate.Length(unsafeParams["user_password"], 8, 100)
		if err != nil {
			return router.BadRequestError(err, "Пароль слишком короткий", "Ваш пароль должен иметь длину не менее 8 символов.")
		}

		ep, err := auth.HashPassword(unsafeParams["user_password"])
		if err != nil {
			return err
		}
		unsafeParams["user_encrypted_password"] = ep
	}

	// Delete password param
	delete(unsafeParams, "user_password")

	return nil
}

/**
* Функция CheckName проверяет имя пользователя.
*
* @param name string Имя пользователя.
* @return error
*/
func CheckName(name string) error {
	n := strings.ToLower(name)

	prohibited := []string{"админ", "путин", "/", "\n", ">", "<", "=", "$", "!", "administrator", "администратор", "&", "*", "{", "}", "admin", "хуй", "пизда"}
	for _, v := range prohibited {
		if strings.Contains(n, v) {
			return fmt.Errorf("Имя содержит запрещённые слова или символы.")
		}
	}

	return nil
}

/**
* Функция GetAvatar получает Gravatar пользователя.
*
* @param email string
* @return string
*/
func GetAvatar(email string) string {
	if email != "" {
		hash := md5.Sum([]byte(email))
		return "https://secure.gravatar.com/avatar/" + hex.EncodeToString(hash[:])
	}

	return ""
}

/**
* Функция ComputeDeltaPower вычисляет дельту для повышения силы пользователя используя логарифмическое распределение.
* Сила - это коэффициент влияния пользователя на других пользователей и их ресурсы.
*
* @return float64 Дельта силы.
*
*/
func (m *User) ComputeDeltaPower() float64 {
	minSize := 0.1
	maxSize := 8.0
	sizeRange := maxSize - minSize

	minCount := math.Log(0 + 1)
	maxCount := math.Log(500 + 1)
	countRange := maxCount - minCount

	newPower := m.Power / 50

	delta := minSize + (math.Log(newPower+1)-minCount)*(sizeRange/countRange)

	return delta
}

/**
* Функция ComputeScoreScore вычисляет дельту для повышения рейтинга пользователя используя логарифмическое распределение.
* Рейтинг - это коэффициент для управления доступом пользователя к ресурсам.
*
* @return float64 Дельта рейтинга.
*
*/
func (m *User) ComputeDeltaScore() float64 {
	return m.ComputeDeltaPower() / 2.73
}

/**
* Функция OwnerBy возвращает true, если ресурс принадлежить пользователю.
*
* @param id int64 Идентификатор пользователя.
* @return bool
*/
func (m *User) OwnedBy(id int64) bool {
	if m.Id == id {
		return true
	}

	return false
}

/**
* Функция Update обновляет пользователя.
*
* @param map[string]string параметры.
* @return error
*/
func (m *User) Update(params map[string]string) error {
	err := validateParams(params)
	if err != nil {
		return err
	}

	// Make sure updated_at is set to the current time
	params["user_updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("user_id=?", m.Id).Update(params)
}

/**
* Destroy this user.
*/
func (m *User) Destroy() error {
	return Query().Where("user_id=?", m.Id).Delete()
}

/**
* URLShow returns the url for this user.
*/
func (m *User) URLShow() string {
	return fmt.Sprintf("/users/%d", m.Id)
}

/**
* Функция проверяет онлайн ли пользователь.
*/
func (m *User) CheckOnline() bool {

	duration := time.Since(m.VisitedAt)
	seconds := duration.Seconds()

	if seconds < 360 {
		return true
	}

	return false
}

/**
* Функция DisplayAvatar возвращает url на avatar.
*
* @param int64 размер аватара
* @return string
*/
func (m *User) DisplayAvatar(s int64) string {
	if !strings.Contains(m.Avatar, "gravatar") {
		return fmt.Sprintf("https://secure.gravatar.com/avatar/?secure=true&d=retro&s=%d", s)
	}

	return m.Avatar + fmt.Sprintf("?secure=true&d=retro&s=%d", s)
}

/**
* Функция DeclensionComments выбирает правильное склонение для колличества комментариев.
*/
func (m *User) DeclensionComments() string {
	switch {
	case (m.CommentCount%10 == 1) && (m.CommentCount%100 != 11):
		return "комментарий"
	case (m.CommentCount%10 >= 2) && (m.CommentCount%10 <= 4) && (m.CommentCount%100 < 10 || m.CommentCount%100 >= 20):
		return "комментария"
	}

	return "комментариев"
}

/**
* Функция DeclensionTopics выбирает правильное склонение для колличества топиков.
*/
func (m *User) DeclensionTopics() string {
	switch {
	case (m.TopicCount%10 == 1) && (m.TopicCount%100 != 11):
		return "топик"
	case (m.TopicCount%10 >= 2) && (m.TopicCount%10 <= 4) && (m.TopicCount%100 < 10 || m.TopicCount%100 >= 20):
		return "топика"
	}

	return "топиков"
}

/**
* Функция CheckCreateCommentLimit проверяет лимит на публикацию комментария.
*
* @return bool
* @return error
*/
func (m *User) CheckCreateCommentLimit() (bool, error) {
	key := fmt.Sprintf("limiter:user:create_comment:%d", m.Id)

	l, err := limiter.New(1, "minute", 2)
	if err != nil {
		return false, err
	}

	limited, _, err := l.GCRA().RateLimit(key, 1)
	if err != nil {
		return false, err
	}

	if !limited {
		return true, nil
	}

	return false, nil
}

/**
* Функция CheckCreateCommentLimit проверяет лимит на публикацию топика.
*
* @return bool
* @return error
*/
func (m *User) CheckCreateTopicLimit() (bool, error) {
	key := fmt.Sprintf("limiter:user:create_topic:%d", m.Id)

	l, err := limiter.New(1, "minute", 1)
	if err != nil {
		return false, err
	}

	limited, _, err := l.GCRA().RateLimit(key, 1)
	if err != nil {
		return false, err
	}

	if !limited {
		return true, nil
	}

	return false, nil
}

/**
* Функция CheckCreateOfferLimit проверяет лимит на публикацию оффера.
*
* @return bool
* @return error
*/
func (m *User) CheckCreateOfferLimit() (bool, error) {
	key := fmt.Sprintf("limiter:user:create_offer:%d", m.Id)

	l, err := limiter.New(1, "minute", 1)
	if err != nil {
		return false, err
	}

	limited, _, err := l.GCRA().RateLimit(key, 1)
	if err != nil {
		return false, err
	}

	if !limited {
		return true, nil
	}

	return false, nil
}


/**
* Функция CheckLoginLimit проверяет лимит на попытки входа.
*
* @return bool
* @return error
*/
func (m *User) CheckLoginLimit() (bool, error) {
	key := fmt.Sprintf("limiter:user:login:%d", m.Id)

	l, err := limiter.New(1, "hour", 5)
	if err != nil {
		return false, err
	}

	limited, _, err := l.GCRA().RateLimit(key, 1)
	if err != nil {
		return false, err
	}

	if !limited {
		return true, nil
	}

	return false, nil
}

/**
* Функция CheckActivationLimit проверяет лимит на попытки входа.
*
* @return bool
* @return int64 колличество попыток.
* @return error
*/
func (m *User) CheckActivationLimit() (bool, int64, error) {
	key := fmt.Sprintf("limiter:user:activation:%d", m.Id)

	l, err := limiter.New(1, "hour", 3)
	if err != nil {
		return false, 0, err
	}

	limited, result, err := l.GCRA().RateLimit(key, 1)
	if err != nil {
		return false, 0, err
	}

	if !limited {
		return true, int64(result.Remaining), nil
	}

	return false, 0, nil
}

/**
 * Функция VerbСonjugation выбирает правильный вариант русского глагола в зависимости от гендера в единственном числе.
 *
 * @param v string глагол, муж.род, ед.число.
*/
func (m *User) VerbСonjugation(v string) string {
	if m.Gender == 0 {
		return v
	}

	if strings.HasSuffix(v, "лся") {
		return strings.Replace(v, "лся", "лась", 1)
	}

	return v + "a"
}

/**
 * Функция SendEmailActivation создает новый токен и отправляет Email для активации профиля пользователя.
 *
 * @return error
*/
func (m *User) SendEmailActivation() error {
	token := auth.BytesToHex(auth.RandomToken())
	if err := m.Update(map[string]string{"user_create_token": token}); err != nil {
		return err
	}

	recipients := []string{m.Email}
	variables := map[string]interface{}{"verification_link": "https://techfront.org/signup?token=" + token, "user_name": m.Name}
	if err := mail.Send(recipients, "Подтверждение регистрации", "component/user/template/mail/verification", variables, nil); err != nil {
		return err
	}

	return nil
}