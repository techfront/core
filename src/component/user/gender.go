package user

import (
	"github.com/techfront/core/src/kernel/view/helper"
)

var (
	GenderM int64 = 0
	GenderW int64 = 1
)

func (m *User) GenderM() bool {
	return m.Gender == GenderM
}

func (m *User) GenderW() bool {
	return m.Gender == GenderW
}

func (m *User) GenderOptions() []helper.ConcreteOption {
	var options []helper.ConcreteOption

	options = append(options, helper.ConcreteOption{Id: GenderM, Name: "М"})
	options = append(options, helper.ConcreteOption{Id: GenderW, Name: "Ж"})

	return options
}

func (m *User) GenderDisplay() string {
	for _, o := range m.GenderOptions() {
		if o.GetId() == m.Gender {
			return o.GetName()
		}
	}
	return "М"
}
