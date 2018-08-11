package topic

import (
	"github.com/fragmenta/query"
	"github.com/techfront/core/src/kernel/view/helper"
)

var (
	FormatTopic    int64 = 0
	FormatNews     int64 = 10
	FormatVideo    int64 = 20
	FormatQuestion int64 = 30
	FormatProject  int64 = 40
	FormatPodcast  int64 = 50
)

/**
* Функция FormatOptions возвращает массив опций для Select.
 */
func (m *Topic) FormatOptions() []helper.ConcreteOption {

	// Пустой список опций
	var options []helper.ConcreteOption

	// Добавление опций
	options = append(options, helper.ConcreteOption{Id: FormatTopic, Name: "Топики"})
	options = append(options, helper.ConcreteOption{Id: FormatNews, Name: "Новости"})
	options = append(options, helper.ConcreteOption{Id: FormatQuestion, Name: "Вопросы"})
	options = append(options, helper.ConcreteOption{Id: FormatVideo, Name: "Видео"})
	options = append(options, helper.ConcreteOption{Id: FormatPodcast, Name: "Подкасты"})
	options = append(options, helper.ConcreteOption{Id: FormatProject, Name: "Проекты"})

	return options
}

/**
* Функция FormatDisplay возвращает название тега.
 */
func (m *Topic) FormatDisplay() string {
	for _, o := range m.FormatOptions() {
		if o.GetId() == m.Id {
			return o.GetName()
		}
	}
	return "Топики"
}

// IsTopic возвращает true если текущий топик - топик
func (m *Topic) IsTopic() bool {
	return m.FormatId == FormatTopic
}

// IsNews возвращает true если текущий топик - новость
func (m *Topic) IsNews() bool {
	return m.FormatId == FormatNews
}

// IsPodcast возвращает true если текущий топик - подкаст
func (m *Topic) IsPodcast() bool {
	return m.FormatId == FormatPodcast
}

// IsVideo возвращает true если текущий топик - видео
func (m *Topic) IsVideo() bool {
	return m.FormatId == FormatVideo
}

// IsQuestion возвращает true если текущий топик - вопрос
func (m *Topic) IsQuestion() bool {
	return m.FormatId == FormatQuestion
}

// IsProject возвращает true если текущий топик - проект
func (m *Topic) IsProject() bool {
	return m.FormatId == FormatProject
}

// Podcast возвращает запрос для получения всех подкастов
func Podcast() *query.Query {
	return Query().Where("topic_id_format=?", FormatPodcast).Order("topic_name asc")
}

// Video возвращает запрос для получения всех видео
func Video() *query.Query {
	return Query().Where("topic_id_format=?", FormatVideo).Order("topic_name asc")
}

// Question возвращает запрос для получения всех вопросов
func Question() *query.Query {
	return Query().Where("topic_id_format=?", FormatQuestion).Order("topic_name asc")
}

// News возвращает запрос для получения всех новостей
func News() *query.Query {
	return Query().Where("topic_id_format=?", FormatNews).Order("topic_name asc")
}

// Project возвращает запрос для получения всех проектов
func Project() *query.Query {
	return Query().Where("topic_id_format=?", FormatProject).Order("topic_name asc")
}
