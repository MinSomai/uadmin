package models

import (
	"fmt"
	"github.com/uadmin/uadmin/interfaces"

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
	interfaces.Model
	ModelName           string `uadmin:"read_only"`
	ModelPK             uint   `uadmin:"read_only"`
	ColumnName          string `uadmin:"read_only"`
	OldValue            string `uadmin:"read_only"`
	NewValue            string
	NewValueDescription string    `uadmin:"read_only"`
	ChangedBy           string    `uadmin:"read_only"`
	ChangeDate          time.Time `uadmin:"read_only"`
	ApprovalAction      ApprovalAction
	ApprovalBy          string     `uadmin:"read_only"`
	ApprovalDate        *time.Time `uadmin:"read_only"`
	ViewRecord          string     `uadmin:"link"`
	UpdatedBy           string     `uadmin:"read_only;hidden;list_exclude"`
}

func (a *Approval) String() string {
	return fmt.Sprintf("Approval for %s.%s %d", a.ModelName, a.ColumnName, a.ModelPK)
}

// Save overides save
func (a *Approval) Save() {
	//if a.ViewRecord == "" {
	//	a.ViewRecord = preloaded.RootURL + a.ModelName + "/" + fmt.Sprint(a.ModelPK)
	//}
	//if modelold.Schema[a.ModelName].FieldByName(a.ColumnName).Type == preloaded.CLIST {
	//	m, _ := modelold.NewModel(a.ModelName, false)
	//	intVal, _ := strconv.ParseInt(a.NewValue, 10, 64)
	//	m.FieldByName(a.ColumnName).SetInt((intVal))
	//	a.NewValueDescription = utils.GetString(m.FieldByName(a.ColumnName).Interface())
	//} else if modelold.Schema[a.ModelName].FieldByName(a.ColumnName).Type == preloaded.CFK {
	//	m, _ := modelold.NewModel(strings.ToLower(modelold.Schema[a.ModelName].FieldByName(a.ColumnName).TypeName), true)
	//	database.Get(m.Interface(), "id = ?", a.NewValue)
	//	a.NewValueDescription = utils.GetString(m.Interface())
	//} else {
	//	a.NewValueDescription = a.NewValue
	//}
	//
	//// Run Approval handle func
	//saveApproval := true
	//if ApprovalHandleFunc != nil {
	//	saveApproval = ApprovalHandleFunc(a)
	//}
	//
	//// Process approval based on the action
	//old := Approval{}
	//if a.ID != 0 {
	//	database.Get(&old, "id = ?", a.ID)
	//}
	//if old.ApprovalAction != a.ApprovalAction {
	//	a.ApprovalBy = a.UpdatedBy
	//	now := time.Now()
	//	a.ApprovalDate = &now
	//	m, _ := modelold.NewModelArray(a.ModelName, true)
	//	model1, _ := modelold.NewModel(a.ModelName, false)
	//	if a.ApprovalAction == a.ApprovalAction.Approved() {
	//		if model1.FieldByName(a.ColumnName).Type().String() == "*time.Time" && a.NewValue == "" {
	//			database.Update(m.Interface(), dialect.GetDB("default").Config.NamingStrategy.ColumnName("", a.ColumnName), nil, "id = ?", a.ModelPK)
	//		} else if modelold.Schema[a.ModelName].FieldByName(a.ColumnName).Type == preloaded.CFK {
	//			database.Update(m.Interface(), dialect.GetDB("default").Config.NamingStrategy.ColumnName("", a.ColumnName)+"_id", a.NewValue, "id = ?", a.ModelPK)
	//		} else {
	//			database.Update(m.Interface(), dialect.GetDB("default").Config.NamingStrategy.ColumnName("", a.ColumnName), a.NewValue, "id = ?", a.ModelPK)
	//		}
	//	} else {
	//		if model1.FieldByName(a.ColumnName).Type().String() == "*time.Time" && a.OldValue == "" {
	//			database.Update(m.Interface(), dialect.GetDB("default").Config.NamingStrategy.ColumnName("", a.ColumnName), nil, "id = ?", a.ModelPK)
	//		} else if modelold.Schema[a.ModelName].FieldByName(a.ColumnName).Type == preloaded.CFK {
	//			database.Update(m.Interface(), dialect.GetDB("default").Config.NamingStrategy.ColumnName("", a.ColumnName)+"_id", a.OldValue, "id = ?", a.ModelPK)
	//		} else {
	//			database.Update(m.Interface(), dialect.GetDB("default").Config.NamingStrategy.ColumnName("", a.ColumnName), a.OldValue, "id = ?", a.ModelPK)
	//		}
	//	}
	//}
	//
	//if !saveApproval {
	//	return
	//}
	//
	//database.Save(a)
}

// ApprovalHandleFunc is a function that could be called during the save process of each approval
var ApprovalHandleFunc func(*Approval) bool
