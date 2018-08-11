package authorise

import (
	"fmt"
	"github.com/fragmenta/server"
	"github.com/techfront/core/src/component/user"
	"github.com/techfront/core/src/kernel/auth"
	"github.com/techfront/core/src/kernel/router"
	"strings"
)

type ResourceModel interface {
	OwnedBy(int64) bool
}

/**
 * Конфигурирование и инициализация
 */
func Setup(s *server.Server) {

	config := s.Configuration()
	auth.HMACKey = auth.HexToBytes(config["hmac_key"])
	auth.SecretKey = auth.HexToBytes(config["secret_key"])
	auth.SessionName = "techfront"

	if s.Production() {
		auth.SecureCookies = true
	}
}

/**
 * Path разрешает путь для текущего пользователя
 */
func Path(context router.Context) error {
	return Resource(context, nil)
}

/**
 * Resource разрешает путь или ресурс для текущего пользователя
 */
func Resource(context router.Context, rm ResourceModel) error {
	userEntity := CurrentUser(context)

	switch userEntity.Role {
	case user.RoleAdmin:
		return authoriseAdmin(context, rm)
	case user.RoleBanned:
		return authoriseBanned(context, rm)
	default:
		return authoriseReader(context, rm)
	}
}

/**
 * ResourceAndAuthenticity разрешает путь или ресурс для текущего пользователя
 */
func ResourceAndAuthenticity(context router.Context, rm ResourceModel) error {
	err := AuthenticityToken(context)
	if err != nil {
		return err
	}

	return Resource(context, rm)
}

/**
 * Возвращает ошибку, если путь или ресурс недоступен для данного типа.
 */
func authoriseAdmin(context router.Context, rm ResourceModel) error {
	return nil
}

/**
 * Возвращает ошибку, если путь или ресурс недоступен для данного типа.
 */
func authoriseReader(context router.Context, rm ResourceModel) error {
	userEntity := CurrentUser(context)

	if context.Path() == "/topics/create" && userEntity.CanSubmit() {
		return nil
	}

	if context.Path() == "/comments/create" && userEntity.CanComment() {
		return nil
	}

	if strings.HasSuffix(context.Path(), "/upvote") && userEntity.CanUpvote() {
		return nil
	}

	if strings.HasSuffix(context.Path(), "/downvote") && userEntity.CanDownvote() {
		return nil
	}

	if rm != nil {
		if rm.OwnedBy(userEntity.Id) {
			return nil
		}
	}

	return fmt.Errorf("Path and Resource not authorized:%s %v", context.Path(), context.Request())

}

/**
 * Возвращает ошибку, если путь или ресурс недоступен для данного типа.
 */
func authoriseBanned(context router.Context, rm ResourceModel) error {
	userEntity := CurrentUser(context)

	if rm != nil {
		if rm.OwnedBy(userEntity.Id) {
			return nil
		}
	}

	return fmt.Errorf("Path and Resource not authorized:%s %v", context.Path(), context.Request())
}
