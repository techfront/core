package user

import (
	"time"
	"github.com/fragmenta/query"
	"github.com/techfront/core/src/kernel/validate"
)

type Favorite struct {
	Id int64
	CreatedAt time.Time
	UserId int64
	TopicId int64
}

/**
* Функция NewFavoriteWithColumns возвращает сущность Favorite.
*
* @param cols map[string]interface{} параметры.
* @return Favorite
*/
func NewFavoriteWithColumns(cols map[string]interface{}) *Favorite {
	f := &Favorite{}
	f.Id = validate.Int(cols["user_favorite_id"])
	f.CreatedAt = validate.Time(cols["user_favorite_created_at"])
	f.UserId = validate.Int(cols["user_favorite_id_user"])
	f.TopicId = validate.Int(cols["user_favorite_id_topic"])

	return f
}

/**
* Функция AddFavorite добавляет закладку пользователя в базу.
*
* @param params map[string]string параметры.
* @return int64 id записи.
* @return error 
*/
func CreateFavorite(params map[string]string) (int64, error) {
	params["user_favorite_created_at"] = query.TimeString(time.Now().UTC())

	return query.New("tf_user_favorite", "user_favorite_id").Insert(params)
}

/**
* Функция DestroyFavorite удаляет закладку пользователя из базы.
*
* @return error
 */
func (m *Favorite) DestroyFavorite() error {
	return query.New("tf_user_favorite", "user_favorite_id").Where("user_favorite_id = ?", m.Id).Delete()
}

/**
* Функция FindFavoriteByTopicId получает закладку пользователя по идентификатору топика.
*
* @param id int64 идентификатор топика.
* @return Favorite, error
*/
func (m *User) FindFavoriteByTopicId(id int64) (*Favorite, error) {
	q := query.New("tf_user_favorite", "user_favorite_id")
	q.Where("user_favorite_id_user = ?", m.Id)
	q.Where("user_favorite_id_topic = ?", id)

	result, err := q.FirstResult()
	if err != nil {
		return nil, err
	}

	return NewFavoriteWithColumns(result), nil
}

/**
* Функция GetFavorites получает закладки пользователя.
*
* @return []Favorite
*/
func (m *User) FindAllFavorites() []*Favorite {
	q := query.New("tf_user_favorite", "user_favorite_id")
	q.Where("user_favorite_id_user = ?", m.Id)
	q.Order("user_favorite_created_at desc, user_favorite_id desc")

	results, err := q.Results()
	if err != nil {
		return nil
	}

	var favorites []*Favorite
	for _, cols := range results {
		favoriteEntity := NewFavoriteWithColumns(cols)
		favorites = append(favorites, favoriteEntity)
	}

	return favorites
}

/**
* Функция IsTopicFavorited проверяет есть ли топик в закладках пользователя.
*
* @param id string идентификатор топика.
* @return bool
*/
func (m *User) IsTopicFavorited(id int64) bool {
	q := query.New("tf_user_favorite", "user_favorite_id")
	q.Where("user_favorite_id_user = ?", m.Id)
	q.Where("user_favorite_id_topic = ?", id)

	result, _ := q.FirstResult()

	if result != nil {
		return true
	}

	return false
}