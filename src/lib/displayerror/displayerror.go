/**
* Библиотека определяет параметры вывода для некоторых типов ошибок.
 */

package displayerror

type displayerror []string

var (
	UnknownError       = displayerror{"Ошибка", "К сожалению, что-то пошло не так."}
	AccessDenied       = displayerror{"Доступ запрещен", "К сожалению, у вас нет доступа к этой странице."}
	AccessDeniedBanned = displayerror{"Доступ запрещен", "К сожалению, вы забанены."}
	AccessDeniedKarma  = displayerror{"Доступ запрещен", "К сожалению, у вас плохая карма."}
	PageNotFound       = displayerror{"Страница не найдена", "К сожалению, страница, которую вы ищите, не может быть найдена."}
	TimeIsOverError    = displayerror{"Ошибка", "К сожалению, время истекло. Повторите попытку."}
)
