package useraction

import (
	"fmt"
	"strings"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/lib/mailchimp"
	"github.com/techfront/core/src/lib/displayerror"
	"github.com/techfront/core/src/kernel/view"
)

func HandleUnsubscribeShow(context router.Context) error {
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	v := view.New("component/user/template/unsubscribe")

	switch params.Get("message") {
	case "error--incorrect_email":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Пожалуйста, проверьте правильность Email."
	case "error--not_found_email":
		v.Vars["message_type"] = "error"
		v.Vars["message"] = "Указанный Email никогда не существовал в базе данных."
	}

	v.Vars["user_email"] = params.Get("email")

	v.Vars["meta_title"] = "Отписка от дайджеста / Техфронт"
	v.Vars["meta_desc"] = "Сообщество энтузиастов, кому по нраву обсуждать технологии, исследования, стартапы и всё то, что происходит в сети."

	return v.Render(context)
}

func HandleUnsubscribe(context router.Context) error {
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err, displayerror.UnknownError...)
	}

	email := params.Get("user_email")
	if len(email) < 3 || !strings.Contains(email, "@") {
		return router.Redirect(context, fmt.Sprintf("/unsubscribe?email=%s&message=error--incorrect_email", email))
	}

	if err := mailchimp.UnsubscribeMember(email); err != nil {
		return router.Redirect(context, fmt.Sprintf("/unsubscribe?email=%s&message=error--not_found_email", email))
	}

	return router.Redirect(context, "/")
}