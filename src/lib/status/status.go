package status

import (
	"github.com/techfront/core/src/kernel/view/helper"
)

// ModelStatus adds status to a model - this may in future be removed and moved into apps as it is frequently modified
type ModelStatus struct {
	Status int64
}

// Status values
// If these need to be substantially modified for a particular model,
// it may be better to move this into the model package concerned and modify as required
var (
	Draft       int64 = 0
	Final       int64 = 10
	Suspended   int64 = 11
	Unavailable int64 = 12
	Rejected    int64 = 13
	Unverified  int64 = 14
	Hidden      int64 = 15
	Published   int64 = 100
	Featured    int64 = 101
)

// StatusOptions returns an array of statuses for a status select
func (m *ModelStatus) StatusOptions() []helper.ConcreteOption {
	var options []helper.ConcreteOption

	options = append(options, helper.ConcreteOption{Id: Draft, Name: "Draft"})
	options = append(options, helper.ConcreteOption{Id: Final, Name: "Final"})
	options = append(options, helper.ConcreteOption{Id: Suspended, Name: "Suspended"})
	options = append(options, helper.ConcreteOption{Id: Unverified, Name: "Unverified"})
	options = append(options, helper.ConcreteOption{Id: Rejected, Name: "Rejected"})
	options = append(options, helper.ConcreteOption{Id: Hidden, Name: "Hidden"})
	options = append(options, helper.ConcreteOption{Id: Published, Name: "Published"})
	options = append(options, helper.ConcreteOption{Id: Featured, Name: "Featured"})

	return options
}

// StatusDisplay returns a string representation of the model status
func (m *ModelStatus) StatusDisplay() string {
	for _, o := range m.StatusOptions() {
		if o.GetId() == m.Status {
			return o.GetName()
		}
	}
	return ""
}

// IsDraft returns true if the status is Draft
func (m *ModelStatus) IsDraft() bool {
	return m.Status == Draft
}

// IsFinal returns true if the status is Final
func (m *ModelStatus) IsFinal() bool {
	return m.Status == Final
}

// IsSuspended returns true if the status is Suspended
func (m *ModelStatus) IsSuspended() bool {
	return m.Status == Suspended
}

// IsUnavailable returns true if the status is unavailable
func (m *ModelStatus) IsUnavailable() bool {
	return m.Status == Unavailable
}

// IsPending returns true if the status is Unverified
func (m *ModelStatus) IsUnverified() bool {
	return m.Status == Unverified
}

// IsRejected returns true if the status is Rejected
func (m *ModelStatus) IsRejected() bool {
	return m.Status == Rejected
}

// IsHidden returns true if the status is HIdden
func (m *ModelStatus) IsHidden() bool {
	return m.Status == Hidden
}

// IsPublished returns true if the status is published *or over*
func (m *ModelStatus) IsPublished() bool {
	return m.Status >= Published // NB >=
}

// IsFeatured returns true if the status is featured
func (m *ModelStatus) IsFeatured() bool {
	return m.Status == Featured
}
