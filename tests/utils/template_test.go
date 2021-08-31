package utils

import (
	"github.com/uadmin/uadmin"
	"testing"
)

type TemplateTestSuite struct {
	uadmin.UadminTestSuite
}

func (s *TemplateTestSuite) TestRenderHTML() {
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTemplate(t *testing.T) {
	uadmin.Run(t, new(TemplateTestSuite))
}
