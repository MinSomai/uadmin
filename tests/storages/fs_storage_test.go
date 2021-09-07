package storages

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sergeyglazyrindev/uadmin"
	"github.com/sergeyglazyrindev/uadmin/core"
	"os"
	"testing"
)

type FsStorageTestSuite struct {
	uadmin.TestSuite
}

func (suite *FsStorageTestSuite) SetupTest() {
	uadmin.NewFullAppForTests()
	err := os.Mkdir(core.CurrentConfig.GetPathToUploadDirectory(), 0755)
	if err != nil {
		assert.True(suite.T(), false, "Couldnt create directory for file uploading")
	}
}

func (suite *FsStorageTestSuite) TearDownSuite() {
	err := os.RemoveAll(core.CurrentConfig.GetPathToUploadDirectory())
	if err != nil {
		assert.True(suite.T(), false, fmt.Errorf("Couldnt remove directory for file uploading"))
	}
	uadmin.ClearTestApp()
}

func (suite *FsStorageTestSuite) TestFullFlow() {
	fsStorage := core.NewFsStorage()
	uploadedFile, _ := fsStorage.Save(&core.FileForStorage{
		Content:           []byte("test"),
		PatternForTheFile: "*.txt",
		Filename:          "uploaded.txt",
	})
	assert.NotEmpty(suite.T(), uploadedFile)
	fileContent, _ := fsStorage.Read(uploadedFile)
	assert.Equal(suite.T(), fileContent, []byte("test"))
	fileStats, _ := fsStorage.Stats(uploadedFile)
	assert.True(suite.T(), fileStats.Size() > 0)
	fileExists, _ := fsStorage.Exists(uploadedFile)
	assert.True(suite.T(), fileExists)
	fileRemoved, _ := fsStorage.Delete(uploadedFile)
	assert.True(suite.T(), fileRemoved)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFsStorage(t *testing.T) {
	uadmin.Run(t, new(FsStorageTestSuite))
}
