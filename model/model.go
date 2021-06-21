package model

import "gorm.io/gorm"

// Model is the standard struct to be embedded
// in any other struct to make it a model for uadmin
type Model struct {
	gorm.Model
}

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
