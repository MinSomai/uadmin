package test

import (
	"github.com/uadmin/uadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UadminTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *UadminTestSuite) SetupTest() {
	db := uadmin.GetDB()
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
	db := uadmin.GetDB()
	db = db.Exec("ROLLBACK")
	if db.Error != nil {
		assert.Equal(suite.T(), true, false, "Couldnt rollback transaction")
	}
}
