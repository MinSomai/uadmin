package models

import (
	"fmt"
	"github.com/uadmin/uadmin/core"

	"time"
)

// ApprovalAction is a selection of approval actions
type ApprovalAction int

// Approved is an accepted change
func (ApprovalAction) Approved() ApprovalAction {
	return 1
}

// Rejected is a rejected change
func (ApprovalAction) Rejected() ApprovalAction {
	return 2
}

func HumanizeApprovalAction(approvalAction ApprovalAction) string {
	switch approvalAction {
	case 1:
		return "approved"
	case 2:
		return "rejected"
	default:
		return "unknown"
	}
}

// Approval is a model that stores approval data
type Approval struct {
	core.Model
	ApprovalAction      ApprovalAction   `uadmin:"list" uadminform:"SelectFieldOptions"`
	ApprovalBy          string           `uadmin:"list" uadminform:"ReadonlyField"`
	ApprovalDate        *time.Time       `uadmin:"list" uadminform:"DatetimeReadonlyFieldOptions"`
	ContentType         core.ContentType `uadmin:"list" uadminform:"ReadonlyField"`
	ContentTypeID       uint
	ModelPK             uint      `uadmin:"list" uadminform:"ReadonlyField" gorm:"default:0"`
	ColumnName          string    `uadmin:"list" uadminform:"ReadonlyField"`
	OldValue            string    `uadmin:"list" uadminform:"ReadonlyField"`
	NewValue            string    `uadmin:"list"`
	NewValueDescription string    `uadmin:"list" uadminform:"ReadonlyField"`
	ChangedBy           string    `uadmin:"list" uadminform:"ReadonlyField"`
	ChangeDate          time.Time `uadmin:"list" uadminform:"DatetimeReadonlyFieldOptions"`
	UpdatedBy           string
}

func (a *Approval) String() string {
	return fmt.Sprintf("Approval for %s.%s %d", a.ContentType.ModelName, a.ColumnName, a.ModelPK)
}
