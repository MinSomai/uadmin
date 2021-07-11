package services

import (
	models2 "github.com/uadmin/uadmin/blueprint/abtest/models"
	"time"
)

func init() {
	go func() {
		//for !database.DbOK {
		//	time.Sleep(time.Second)
		//}
		go abTestService()
	}()
}

func abTestService() {
	for {
		if models2.AbTestCount != 0 {
			models2.SyncABTests()
		}
		time.Sleep(time.Second * 10)
	}
}
