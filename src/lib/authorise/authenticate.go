package authorise

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/techfront/core/src/kernel/auth"
	"github.com/techfront/core/src/kernel/router"
	"github.com/techfront/core/src/component/user"
)

func CurrentUser(context router.Context) *user.User {

	userEntity := &user.User{}

	session, err := auth.Session(context.Writer(), context.Request())
	if err != nil {
		return userEntity
	}

	var id int64
	ids := session.Get(auth.SessionUserKey)
	if len(ids) > 0 {
		id, err = strconv.ParseInt(ids, 10, 64)
		if err != nil {
			return userEntity
		}
	}

	if id != 0 {
		if context.Get("current_user") != nil {
			return context.Get("current_user").(*user.User)
		}

		u, err := user.Find(id)
		if err != nil {
			return userEntity
		}
		userEntity = u
	}

	return userEntity
}

func AuthenticityToken(context router.Context) error {
	if context.Request().Method == http.MethodGet {
		return nil
	}

	token := fmt.Sprintf("%v", context.Get(auth.SessionTokenKey))

	err := auth.CheckAuthenticityToken(token, context.Request())
	if err != nil {
		session, err := auth.SessionGet(context.Request())
		if err != nil {
			return err
		}

		session.Clear(context.Writer())
	}

	return nil
}
