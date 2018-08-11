package offer

import (
	"github.com/fragmenta/query"
	"github.com/techfront/core/src/kernel/view/helper"
)

var (
	FormatOffer    int64 = 0
	FormatJob	int64 = 10
)

/**
* Функция FormatOptions возвращает массив опций для Select.
 */
func (m *Offer) FormatOptions() []helper.ConcreteOption {

	// Пустой список опций
	var options []helper.ConcreteOption

	// Добавление опций
	options = append(options, helper.ConcreteOption{Id: FormatOffer, Name: "Офферы"})
	options = append(options, helper.ConcreteOption{Id: FormatJob, Name: "Вакансии"})

	return options
}

/**
* Функция FormatDisplay возвращает название тега.
 */
func (m *Offer) FormatDisplay() string {
	for _, o := range m.FormatOptions() {
		if o.GetId() == m.Id {
			return o.GetName()
		}
	}
	return "Офферы"
}

// IsOffer возвращает true если текущий оффер - оффер
func (m *Offer) IsOffer() bool {
	return m.FormatId == FormatOffer
}

// IsNews возвращает true если текущий оффер - вакансия
func (m *Offer) IsJob() bool {
	return m.FormatId == FormatJob
}

// News возвращает запрос для получения всех вакансий
func Jobs() *query.Query {
	return Query().Where("offer_id_format=?", FormatJob).Order("offer_name asc")
}
