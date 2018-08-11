package topic

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/techfront/core/src/kernel/model"
	"github.com/techfront/core/src/kernel/validate"
	"github.com/fragmenta/query"

	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/cache"
	"github.com/techfront/core/src/lib/status"
)

type Topic struct {
	model.Model
	status.ModelStatus
	Name string
	Url string
	Text string
	Thumbnail string
	FormatId int64
	UserId int64
	UpvoteCount int64
	DownvoteCount int64
	FlagCount int64
	CommentCount int64
	Rank float64
	Score float64
	Points int64
	UserData userData
}

type userData struct {
	Id int64
	Name string
	Avatar string
	Gender int64
}

/**
* AllowedParams возвращает массив разрешенных параметров.
*/
func AllowedParams() []string {
	return []string{"topic_name", "topic_count_upvote", "topic_text", "topic_url", "topic_id_format", "topic_thumbnail"}
}

/**
* AllowedParamsAdmin возвращает массив разрешенных параметров для Администратора.
*/
func AllowedParamsAdmin() []string {
	return []string{"topic_name", "topic_status", "topic_id_user", "topic_thumbnail", "topic_text", "topic_url", "topic_id_format", "topic_count_comment", "topic_count_upvote", "topic_count_downvote", "topic_count_flag", "topic_tw_posted_at", "topic_vk_posted_at", "topic_newsletter_at"}
}

/**
* AllowedUpdateParams возвращает массив разрешенных параметров при обновлении.
*/
func AllowedUpdateParams() []string {
	return []string{"topic_name", "topic_text", "topic_id_format"}
}

/**
* NewWithColumns инициализирует сущность топика с данными.
*/
func NewWithColumns(cols map[string]interface{}) *Topic {
	t := New()
	t.Id = validate.Int(cols["topic_id"])
	t.CreatedAt = validate.Time(cols["topic_created_at"])
	t.UpdatedAt = validate.Time(cols["topic_updated_at"])
	t.Status = validate.Int(cols["topic_status"])
	t.Name = validate.String(cols["topic_name"])
	t.Url = validate.String(cols["topic_url"])
	t.Text = validate.String(cols["topic_text"])
	t.Thumbnail = validate.String(cols["topic_thumbnail"])

	t.FormatId = validate.Int(cols["topic_id_format"])
	t.UserId = validate.Int(cols["topic_id_user"])

	t.UpvoteCount = validate.Int(cols["topic_count_upvote"])
	t.DownvoteCount = validate.Int(cols["topic_count_downvote"])
	t.FlagCount = validate.Int(cols["topic_count_flag"])
	t.CommentCount = validate.Int(cols["topic_count_comment"])

	// Вычисление Score
	t.Score = t.ComputeScore()

	// Вычисление Rank
	t.Rank = t.ComputeRank()

	// Вычисление Points
	t.Points = t.ComputePoints()

	// Получение данных автора
	t.UserData = t.GetUserData()

	return t
}

/**
* Функция New инициализирует сущность топика.
*/
func New() *Topic {
	t := &Topic{}
	t.Model.Init()
	t.Status = status.Published
	t.TableName = "tf_topic"
	t.KeyName = "topic_id"
	return t
}

/**
* Функция Create создаёт запись в базе данных с параметрами.
*
* @param params map[string]string Параметры.
* @return int64, error ID созданного топика или ошибка.
*/
func Create(params map[string]string) (int64, error) {

	// Проверка параметров
	err := validateParams(params)
	if err != nil {
		return 0, err
	}

	// Получение времени создания/обновления
	params["topic_created_at"] = query.TimeString(time.Now().UTC())
	params["topic_updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

/**
* validateParams проверяет переданные параметры.
*
* @param params map[string]string Параметры.
* @return error Ошибка, если есть, иначе - nil.
*/
func validateParams(params map[string]string) error {

	if len(params["topic_name"]) > 0 {
		err := validate.Length(params["topic_name"], 10, 1000)
		if err != nil {
			return router.BadRequestError(err, "Неправильный заголовок", "Заголовок должен содержать не менее 10, но не более 150 символов.")
		}
	}

	if len(params["topic_url"]) > 0 {
		err := validate.Length(params["topic_url"], 5, 1000)
		if err != nil {
			return router.BadRequestError(err, "Неправильная ссылка", "Ссылка должна содержать более пяти символов.")
		}

		if !strings.HasPrefix(params["topic_url"], "http://") && !strings.HasPrefix(params["topic_url"], "https://") {
			return router.BadRequestError(nil, "Неправильная ссылка", "Ссылка должна начинаться с https:// или http://")
		}
	}

	return nil
}

/**
* Функция Find возвращает одну запись по id.
*/
func Find(id int64) (*Topic, error) {
	result, err := Query().Where("topic_id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

/**
* Функция FindAll возвращает все результаты по запросу.
*/
func FindAll(q *query.Query) ([]*Topic, error) {

	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	var topics []*Topic
	for _, cols := range results {
		p := NewWithColumns(cols)
		topics = append(topics, p)
	}

	return topics, nil
}

/**
* Функция Query возвращает инициализированный запрос для топиков.
*/
func Query() *query.Query {
	p := New()
	return query.New(p.TableName, p.KeyName)
}

/**
* Where returns a Where query for topics with the arguments supplied.
*/
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

/**
* Published returns a query for all topics with status >= published.
*/
func Published() *query.Query {
	return Query().Where("topic_status >= ?", status.Published)
}

/**
* Popular returns a query for all topics with score over a certain threshold.
*/
func Popular() *query.Query {
	return Query().Where("score(topic_count_upvote, topic_count_downvote, topic_count_flag) > 0")
}

/**
* Функция Favorited возвращает запрос для получения топиков из закладок.
*
* @param id int64 идентификатор пользователя.
* @return *query.Query
*/
func Favorited(id int64) *query.Query {
	q := Query().Select("SELECT tf_topic.* FROM tf_topic INNER JOIN tf_user_favorite ON tf_topic.topic_id = tf_user_favorite.user_favorite_id_topic")
	q.Where("tf_user_favorite.user_favorite_id_user = ?", id)
	q.Order("tf_user_favorite.user_favorite_created_at desc")

	return q
}

/**
* Функция GetNewestCount получает колличество опубликованных топиков за определенное время.
*
* @param t string Промежуток времени (например "1h" или "3600s" равняется одному часу).
* @return result int64
*
* @optimization Глубина поиска - 30 записей, кэширование результата.
*/
func GetNewestCount(t string) int64 {
	key := fmt.Sprintf("cache:topic:newest_count_%s", t)
	var result int64
	if err := cache.Get(key, &result); err == nil {
		return result
	}

	duration, err := time.ParseDuration(t)
	if err != nil {
		return 0
	}

	moment := time.Now().Add(-duration)
	count, err := Published().Where("topic_created_at >= ?", moment).Limit(30).Count()
	if err != nil {
		return 0
	}

	result = count
	if err := cache.Set(key, result, 360); err != nil {
		return 0
	}

	return result
}

/**
* Функция GetUserData() получает и задает данные о пользователе, в зависимости от топика.
*
* @optiomization Интегрированно кеширование.
*/
func (m *Topic) GetUserData() userData {
	key := fmt.Sprintf("cache:topic:user_data_%d", m.UserId)
	var result userData
	if err := cache.Get(key, &result); err == nil {
		return result
	}

	userEntity, err := user.Find(m.UserId)
	if err != nil {
		return result
	}

	result.Id = userEntity.Id
	result.Name = userEntity.Name
	result.Avatar = userEntity.Avatar
	result.Gender = userEntity.Gender

	if err := cache.Set(key, result, 360); err != nil {
		return result
	}

	return result
}

/**
 * Функция ComputeRank Вычисляет актуальность топика.
 *
 * @return float64
*/
func (m *Topic) ComputeRank() float64 {
	duration := time.Since(m.CreatedAt)
	hours := duration.Hours()
	gravity := 1.6

	return m.Score / math.Pow(hours+2, gravity)
}

/**
 * Функция ComputeScore Вычисляет рейтинг топика по Вильсону.
 *
 * @return float64
*/
func (m *Topic) ComputeScore() float64 {
	z := 1.96
	n := float64(m.UpvoteCount + m.DownvoteCount + (m.FlagCount * 5))
	phat := float64(m.UpvoteCount) / n

	return (phat + z*z/(2*n) - z*math.Sqrt((phat*(1-phat)+z*z/(4*n))/n)) / (1 + z*z/n)
}

/**
 * Функция СomputePoints Вычисляет разность между положительными и отрицательными отметками.
 *
 * @return int64
*/
func (m *Topic) ComputePoints() int64 {
	return m.UpvoteCount - m.DownvoteCount - (m.FlagCount * 5)
}

/**
 * Функция OwnedBy возвращает true, если ресурс пренадлежит пользователю.
*/
func (m *Topic) OwnedBy(id int64) bool {
	if m.UserId == id {
		return true
	}

	return false
}

/**
* Функция Related возвращает запрос для получения похожих топиков.
*
* @return *query.Query
*/
func (m *Topic) Related() *query.Query {
	wildcard := "%" + m.DisplayDomain() + "%"

	return Published().Where("(topic_url ILIKE ? OR topic_id_format = ?) AND topic_id != ?", wildcard, m.FormatId, m.Id)
}

/**
 * Update sets the record in the database from params.
*/
func (m *Topic) Update(params map[string]string) error {

	// Проверка параметров
	err := validateParams(params)
	if err != nil {
		return err
	}

	// Получение времени обновления
	params["topic_updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("topic_id=?", m.Id).Update(params)
}

/**
* Функция Destroy удаляет запись из базы данных.
*/
func (m *Topic) Destroy() error {
	return Query().Where("topic_id=?", m.Id).Delete()
}

/**
* URLShow returns the url for this topic.
*/
func (m *Topic) URLShow() string {
	return fmt.Sprintf("/topics/%d", m.Id)
}

/**
* Функция возвращает полную ссылку к топику.
*/
func (m *Topic) FullURLShow() string {
	return fmt.Sprintf("https://techfront.org/topics/%d", m.Id)
}

/**
* DestinationURL returns the URL of the topic.
*/
func (m *Topic) DestinationURL() string {
	if !m.IsVideo() && !m.IsQuestion() && m.Url != "" && len(m.Text) <= 0 {
		return fmt.Sprintf("https://techfront.org/away?to=%s", m.Url)
	}

	return m.URLShow()
}

/**
* RSSDestinationURL returns the URL of the topic.
*/
func (m *Topic) RSSDestinationURL() string {
	if !m.IsVideo() && !m.IsQuestion() && m.Url != "" && len(m.Text) <= 0 {
		return fmt.Sprintf("https://techfront.org/away?to=%s", m.Url)
	}

	return m.FullURLShow()
}

/**
* Функция IsYouTube возвращает true, если ссылка на YouTube.
*/
func (m *Topic) IsYouTube() bool {
	return strings.Contains(m.Url, "youtube.com/watch?v=")
}

/**
* Функция YouTubeURL приводит youtube url к каноническому виду.
*/
func (m *Topic) YouTubeURL() string {
	url := strings.Replace(m.Url, "https://m.youtube.com", "https://www.youtube.com", 1)
	url = strings.Replace(url, "watch?v=", "embed/", 1)
	return url
}

/**
* DisplayNegativePoints returns a negative score score or 0 if score is above 0.
*/
func (m *Topic) DisplayNegativePoints() int64 {
	if m.Points > 0 {
		return 0
	}
	return m.Points
}

/**
* Функция DisplayCommentCount возвращает колличество комментариев.
*/
func (m *Topic) DisplayCommentCount() string {
	if m.CommentCount > 0 {
		return fmt.Sprintf("%d", m.CommentCount)
	}
	return "…"
}

/**
* Domain returns the domain of the topic URL.
*/
func (m *Topic) DisplayDomain() string {
	parts := strings.Split(m.Url, "/")
	if len(parts) > 2 {
		return strings.Replace(parts[2], "www.", "", 1)
	}

	if len(m.Url) > 0 {
		return m.Url
	}

	return "techfront.org"
}

/**
* Функция DisplayAvatar возвращает url на avatar.
*
* @param int64 Размер аватара
* @return string
*/
func (m *Topic) DisplayAvatar(s int64) string {
	if !strings.Contains(m.UserData.Avatar, "gravatar") {
		return fmt.Sprintf("https://secure.gravatar.com/avatar/?secure=true&d=retro&s=%d", s)
	}

	return m.UserData.Avatar + fmt.Sprintf("?secure=true&d=retro&s=%d", s)
}

/**
* Функция DisplayAction отображает действие пользователя.
* 
* @return string
*/
func (m *Topic) DisplayAction() string {
	var action string

	if m.IsQuestion() {
		if m.UserData.Gender == 0 {
			action = "задал вопрос"
		} else {
			action = "задала вопрос"
		}
	} else {
		action = m.VerbСonjugation("поделился")
	}

	return action
}

/**
* Функция DeclensionComments выбирает правильное склонение для колличества комментариев.
*/
func (m *Topic) DeclensionComments() string {
	switch {
	case (m.CommentCount%10 == 1) && (m.CommentCount%100 != 11):
		return "комментарий"
	case (m.CommentCount%10 >= 2) && (m.CommentCount%10 <= 4) && (m.CommentCount%100 < 10 || m.CommentCount%100 >= 20):
		return "комментария"
	}

	return "комментариев"
}

/**
* Функция DeclensionComments выбирает правильное склонение для колличества голосов.
*/
func (m *Topic) DeclensionPoints() string {
	switch {
	case (m.Points%10 == 1) && (m.Points%100 != 11):
		return "голос"
	case (m.Points%10 >= 2) && (m.Points%10 <= 4) && (m.Points%100 < 10 || m.Points%100 >= 20):
		return "голоса"
	}

	return "голосов"
}

/**
* Функция VerbСonjugation выбирает правильную форму русского глагола в зависимости от гендера в единственном числе.
*
* v string глагол, муж.род, ед.число.
*/
func (m *Topic) VerbСonjugation(v string) string {
	if m.UserData.Gender == 0 {
		return v
	}

	if strings.HasSuffix(v, "лся") {
		return strings.Replace(v, "лся", "лась", 1)
	}

	return v + "a"
}
