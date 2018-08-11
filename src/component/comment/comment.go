package comment

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/techfront/core/src/kernel/model"
	"github.com/techfront/core/src/kernel/validate"
	"github.com/fragmenta/query"

	"github.com/techfront/core/src/lib/cache"
	"github.com/techfront/core/src/lib/status"

	"github.com/techfront/core/src/component/topic"
	"github.com/techfront/core/src/component/user"
)

const RANK_ORDER = "score(comment_count_upvote, comment_count_downvote, comment_count_flag) desc, comment_id desc"

/**
* Comment handles saving and retreiving comments from the database.
 */
type Comment struct {
	model.Model
	status.ModelStatus
	Text string
	UserId int64
	TopicId int64
	ParentId int64
	DottedIds string
	Children []*Comment
	UpvoteCount int64
	DownvoteCount int64
	FlagCount int64
	Rank float64
	Score float64
	Points int64
	UserData userData
	TopicData topicData
}

type userData struct {
	Id int64
	Name string
	Avatar string
	Gender int64
}

type topicData struct {
	Id int64
	Name string
	CommentCount int64
}

/**
* AllowedParams returns an array of allowed param keys.
*/
func AllowedParams() []string {
	return []string{"comment_text", "comment_id_parent", "comment_id_topic"}
}

/**
* AllowedParamsAdmin returns an array of allowed param keys.
*/
func AllowedParamsAdmin() []string {
	return []string{"comment_status", "comment_id_user", "comment_id_parent", "comment_count_upvote", "comment_count_downvote", "comment_count_flag", "comment_id_topic", "comment_text", "comment_dotted_ids"}
}

/**
* NewWithColumns creates a new comment instance and fills it with data from the database cols provided.
 */
func NewWithColumns(cols map[string]interface{}) *Comment {
	c := New()
	c.Id = validate.Int(cols["comment_id"])
	c.CreatedAt = validate.Time(cols["comment_created_at"])
	c.UpdatedAt = validate.Time(cols["comment_updated_at"])
	c.Text = validate.String(cols["comment_text"])
	c.Status = validate.Int(cols["comment_status"])
	c.UserId = validate.Int(cols["comment_id_user"])
	c.TopicId = validate.Int(cols["comment_id_topic"])
	c.ParentId = validate.Int(cols["comment_id_parent"])
	c.DottedIds = validate.String(cols["comment_dotted_ids"])
	c.UpvoteCount = validate.Int(cols["comment_count_upvote"])
	c.DownvoteCount = validate.Int(cols["comment_count_downvote"])
	c.FlagCount = validate.Int(cols["comment_count_flag"])

	// Вычисление Score
	c.Score = c.ComputeScore()

	// Вычисление Rank
	c.Rank = c.ComputeRank()

	// Вычисление Points
	c.Points = c.ComputePoints()

	// Получение данных топика
	c.TopicData = c.GetTopicData()

	// Получение данных пользователя
	c.UserData = c.GetUserData()

	return c
}

/**
* New creates and initialises a new comment instance.
*/
func New() *Comment {
	c := &Comment{}
	c.Model.Init()
	c.Status = status.Published
	c.TableName = "tf_comment"
	c.KeyName = "comment_id"
	return c
}

/**
* Create inserts a new record in the database using params, and returns the newly created id.
*/
func Create(params map[string]string) (int64, error) {
	err := validateParams(params)
	if err != nil {
		return 0, err
	}

	params["comment_created_at"] = query.TimeString(time.Now().UTC())
	params["comment_updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

/**
* validateParams checks these params pass validation checks.
*/
func validateParams(params map[string]string) error {
	err := validate.Length(params["comment_id"], 0, -1)
	if err != nil {
		return err
	}

	return err
}

/**
* Find returns a single record by id in params.
*/
func Find(id int64) (*Comment, error) {
	result, err := Query().Where("comment_id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

/**
* FindAll returns all results for this query.
*/
func FindAll(q *query.Query) ([]*Comment, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of topics constructed from the results
	var comments []*Comment
	for _, cols := range results {
		c := NewWithColumns(cols)
		comments = append(comments, c)
	}

	return comments, nil
}

/**
* FindAll returns all results for this query
*/
func FindAllWithChild(q *query.Query) ([]*Comment, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Construct an array of comments constructed from the results
	// We do things a little differently, as we have a tree of comments
	// root comments are added to the list, others are held in another list
	// and added as children to rootComments

	var rootComments, childComments []*Comment
	for _, cols := range results {
		c := NewWithColumns(cols)
		if c.Root() {
			rootComments = append(rootComments, c)
		} else {
			childComments = append(childComments, c)
		}
	}

	// Now walk through child comments, assigning them to their parent

	// Walk through comments, adding those with no parent id to comments list
	// and others to the parent comment in root comments
	for _, c := range childComments {
		found := false
		for _, p := range rootComments {
			if p.Id == c.ParentId {
				p.Children = append(p.Children, c)
				found = true
				break
			}
		}
		if !found {
			for _, p := range childComments {
				if p.Id == c.ParentId {
					p.Children = append(p.Children, c)
					break
				}
			}
		}
	}

	return rootComments, nil
}

/**
* Query returns a new query for comments.
*/
func Query() *query.Query {
	p := New()
	return query.New(p.TableName, p.KeyName)
}

/**
* Published returns a query for all comments with status >= published.
*/
func Published() *query.Query {
	return Query().Where("comment_status>=?", status.Published)
}

/**
* Where returns a Where query for comments with the arguments supplied.
*/
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

/**
* Функция SetUserData получает и задает данные о пользователе, в зависимости от комментария.
*
* @optimization Интегрированно кеширование.
*/
func (m *Comment) GetUserData() userData {
	key := fmt.Sprintf("cache:comment:user_data_%d", m.UserId)
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
* GetTopicData() получает и задает данные о топике, в зависимости от комментария.
*
* @optimization Интегрированно кеширование.
*/
func (m *Comment) GetTopicData() topicData {
	key := fmt.Sprintf("cache:comment:topic_data_%d", m.TopicId)
	var result topicData
	if err := cache.Get(key, &result); err == nil {
		return result
	}

	topicEntity, err := topic.Find(m.TopicId)
	if err != nil {
		return result
	}

	result.Id = topicEntity.Id
	result.Name = topicEntity.Name
	result.CommentCount = topicEntity.CommentCount
	if err := cache.Set(key, result, 360); err != nil {
		return result
	}

	return result
}

/**
* GetRecentComments получает последние комментарии.
*
* @optimization кеширование результата для заданного колличества.
*/
func GetRecentComments(n int) []*Comment {
	key := fmt.Sprintf("cache:comment:recent_%d", n)
	var results []*Comment
	if err := cache.Get(key, &results); err == nil {
		return results
	}

	commentQuery := Query().Limit(n).Where("comment_status != ?", 15).Order("comment_created_at desc")
	commentList, err := FindAllWithChild(commentQuery)
	if err != nil {
		return results
	}

	results = commentList
	if err := cache.Set(key, results, 360); err != nil {
		return results
	}

	return results
}

/**
* Функция проверяет существование комментария.
*/
func (m *Comment) CheckExist() bool {
	if m.UserData.Name != "" && m.TopicData.Name != "" && m.Status != 15 {
		return true
	}

	return false
}

/**
* Функция СomputeRank Вычисляет актуальность топика.
*
* @return float64
*/
func (m *Comment) ComputeRank() float64 {
	duration := time.Since(m.CreatedAt)
	hours := duration.Hours()
	gravity := 1.6

	return m.Score / math.Pow(hours+2, gravity)
}

/**
* Функция ComputeScore Вычисляет рейтинг коммента по Вильсону.
*
* @return float64
*/
func (m *Comment) ComputeScore() float64 {
	// Точность ~ 95%
	z := 1.64485

	// Сумма всех отметок
	n := float64(m.UpvoteCount + m.DownvoteCount + (m.FlagCount * 5))

	// Доля положительных отметок
	phat := float64(m.UpvoteCount) / n

	return (phat + z*z/(2*n) - z*math.Sqrt((phat*(1-phat)+z*z/(4*n))/n)) / (1 + z*z/n)
}

/**
* Функция ComputePoints Вычисляет разность между положительными и отрицательными отметками.
*
* @return float64
*/
func (m *Comment) ComputePoints() int64 {
	return m.UpvoteCount - m.DownvoteCount - (m.FlagCount * 5)
}

/**
* Update sets the record in the database from params.
 */
func (m *Comment) Update(params map[string]string) error {
	err := validateParams(params)
	if err != nil {
		return err
	}

	params["comment_updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("comment_id=?", m.Id).Update(params)
}

/**
* Функция удаляет комментарий.
*/
func (m *Comment) Destroy() error {
	userEntity, err := user.Find(m.UserId)
	if err != nil {
		return err
	}

	topicEntity, err := topic.Find(m.TopicId)
	if err != nil {
		return err
	}

	topicParams := map[string]string{"topic_count_comment": fmt.Sprintf("%d", topicEntity.CommentCount-1)}
	err = topicEntity.Update(topicParams)
	if err != nil {
		return err
	}

	userParams := map[string]string{"user_count_comment": fmt.Sprintf("%d", userEntity.CommentCount-1)}
	err = userEntity.Update(userParams)
	if err != nil {
		return err
	}

	return Query().Where("comment_id=?", m.Id).Delete()
}

/**
* URLTopic returns the internal resource URL for our topic.
*/
func (m *Comment) URLTopic() string {
	return fmt.Sprintf("/topics/%d", m.TopicId)
}

/**
* DisplayNegativePoints returns a negative score between 0 and 5 (positive score return 0, below -6 returns 6).
*/
func (m *Comment) DisplayNegativePoints() int64 {
	if m.Points > 0 {
		return 0
	}
	if m.Points < -6 {
		return 6
	}

	return m.Points
}

/**
* Функция DisplayAvatar возвращает url на avatar.
*
* @param int64 Размер аватара
* @return string
*/
func (m *Comment) DisplayAvatar(s int64) string {
	if !strings.Contains(m.UserData.Avatar, "gravatar") {
		return fmt.Sprintf("https://secure.gravatar.com/avatar/?secure=true&d=retro&s=%d", s)
	}

	return m.UserData.Avatar + fmt.Sprintf("?secure=true&d=retro&s=%d", s)
}

/**
* Level returns the nesting level of this comment, based on dotted_ids.
*/
func (m *Comment) Level() int64 {
	if m.ParentId > 0 {
		return int64(strings.Count(m.DottedIds, "."))
	}

	return 0
}

/**
* Root returns true if this is a root comment.
*/
func (m *Comment) Root() bool {
	return m.ParentId == 0
}

/**
* Функция VerbСonjugation выбирает правильный вариант русского глагола в зависимости от гендера в единственном числе.
*
* @param v string глагол, муж.род (извинения феменисткам), ед.число.
*/
func (m *Comment) VerbСonjugation(v string) string {
	if m.UserData.Gender == 0 {
		return v
	}

	if strings.HasSuffix(v, "лся") {
		return strings.Replace(v, "лся", "лась", 1)
	}

	return v + "a"
}
