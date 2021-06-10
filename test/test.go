package test

import (
	"github.com/uadmin/uadmin/dialect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UadminTestSuite struct {
	suite.Suite
}

func (suite *UadminTestSuite) SetupTest() {
	db := dialect.GetDB()
	db = db.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")
	if db.Error != nil {
		assert.Equal(suite.T(), true, false, "Couldnt setup isolation level for db")
	}
	db = db.Exec("BEGIN")
	if db.Error != nil {
		assert.Equal(suite.T(), true, false, "Couldnt start transaction")
	}
}

func (suite *UadminTestSuite) TearDownSuite() {
	db := dialect.GetDB()
	db = db.Exec("ROLLBACK")
	if db.Error != nil {
		assert.Equal(suite.T(), true, false, "Couldnt rollback transaction")
	}
}
