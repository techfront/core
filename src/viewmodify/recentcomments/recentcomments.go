package recentcomments

import (
	"github.com/techfront/core/src/component/comment"
	"github.com/techfront/core/src/kernel/view"
)

/**
 * Модификатор шаблонов, добавляет данные для виджета "Сейчас обсуждают"
 *
 */
func Modify(context view.Context, v *view.View) {
	v.Vars["comment_list_recent"] = comment.GetRecentComments(10)
}
