package models

import "github.com/sergeyglazyrindev/uadmin/core"

type Todo struct {
	core.Model
	TaskAlias string
	TaskDescription string
}
