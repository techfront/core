package useraction

import (
	"strings"

	"github.com/techfront/core/src/kernel/response"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/mailchimp"
)

func HandleSubscribeEndpoint(context router.Context) error {

	w := context.Writer()

	params, err := context.Params()
	if err != nil {
		return response.Send(w, 500, response.InternalError("Ошибка 500. Что-то пошло не так."))
	}

	email := params.Get("email")
	if len(email) > 0 {
		if len(email) < 3 || !strings.Contains(email, "@") {
			return response.Send(w, 422, response.Unprocessable("Ошибка 422. Вероятно, вы допустили ошибку."))
		}
	} else {
		return response.Send(w, 422, response.Unprocessable("Ошибка 422. Поле \"Email\" не заполнено."))
	}

	if err := mailchimp.AddMember(email); err != nil {
		return response.Send(w, 400, response.BadRequest("Ошибка 400. Вероятно, вы уже были подписаны."))
	}

	return response.Send(w, 200, &response.Response{
		Data: nil,
		Meta: map[string]interface{}{
			"_notice": "Спасибо, вы успешно подписались.",
		},
	})
}
