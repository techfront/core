package router

var (
	MessageDefault int64 = 0
	MessageError   int64 = 1
	MessageSuccess int64 = 2
	MessageWarning int64 = 3
)

/**
* Тип Message это структура сообщения в контексте.
* Как правило, сообщение может быть добавлено через хендлер или посредник (middleware),
* Затем передано в шаблонизатор через модификатор (viewmodifier).
 */
type Message struct {
	uid string
	style int64
	body  string
}

/**
* Функция Uid возвращает уникальный идентификатор сообщения.
*/
func (m Message) Uid() string {
	return m.uid
}

/**
* Функция Style возвращает тип(стиль) сообщения.
*
* @return int64
 */
func (m Message) Style() int64 {
	return m.style
}

/**
* Функция Body возвращает тело сообщения.
*
* @return string
 */
func (m Message) Body() string {
	return m.body
}
