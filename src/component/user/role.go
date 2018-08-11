package user

import (
	"github.com/fragmenta/query"
	"github.com/techfront/core/src/kernel/view/helper"
)

var (
	RoleAnon   int64 = 0
	RoleBanned int64 = 10
	RoleReader int64 = 20
	RoleEditor int64 = 50
	RoleAdmin  int64 = 100
)

func (m *User) Anon() bool {
	return m.Role == RoleAnon
}

func (m *User) Banned() bool {
	return m.Role == RoleBanned
}

func (m *User) Reader() bool {
	return m.Role == RoleReader
}

func (m *User) Editor() bool {
	return m.Role == RoleEditor
}

func (m *User) Admin() bool {
	return m.Role == RoleAdmin
}

func (m *User) RoleOptions() []helper.ConcreteOption {
	var options []helper.ConcreteOption

	options = append(options, helper.ConcreteOption{Id: RoleBanned, Name: "Заблокированный"})
	options = append(options, helper.ConcreteOption{Id: RoleReader, Name: "Зарегистрированный"})
	options = append(options, helper.ConcreteOption{Id: RoleEditor, Name: "Редактор"})
	options = append(options, helper.ConcreteOption{Id: RoleAdmin, Name: "Администратор"})

	return options
}

func (m *User) RoleDisplay() string {
	for _, o := range m.RoleOptions() {
		if o.GetId() == m.Role {
			return o.GetName()
		}
	}
	return "Зарегистрированный"
}

func Admins() *query.Query {
	return Query().Where("user_role=?", RoleAdmin).Order("user_name asc")
}

func Editors() *query.Query {
	return Query().Where("user_role=?", RoleEditor).Order("user_name asc")
}

func Readers() *query.Query {
	return Query().Where("user_role=?", RoleReader).Order("user_name asc")
}

func Banneds() *query.Query {
	return Query().Where("user_role=?", RoleBanned).Order("user_name asc")
}
