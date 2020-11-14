package adapter_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"testing"
)

var db2Name = fmt.Sprintf("test_%d.db", rand.Int())

type Post2 struct {
	ID     uint `gorm:"primary_key"`
	Number int
}

type GORM2AdapterTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *GORM2AdapterTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(db2Name), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("setup test: %s", err))
	}

	suite.db = db
	suite.db.AutoMigrate(&Post2{})

	for i := 1; i <= 100; i++ {
		p := Post2{
			Number: i,
		}

		suite.db.Save(&p)
	}
}

func (suite *GORM2AdapterTestSuite) TearDownTest() {
	sqlDB, err := suite.db.DB()
	if err != nil {
		panic(fmt.Errorf("tear down test: %s", err))
	}
	if err := sqlDB.Close(); err != nil {
		panic(fmt.Errorf("tear down test: %s", err))
	}
	if err := os.Remove(db2Name); err != nil {
		panic(fmt.Errorf("tear down test: %s", err))
	}
}

func (suite *GORM2AdapterTestSuite) TestFirstPage() {
	q := suite.db.Model(Post2{})
	p := paginator.New(adapter.NewGORM2Adapter(q), 10)

	assert.Equal(suite.T(), 10, p.PageNums())
	assert.Equal(suite.T(), 1, p.Page())
	assert.True(suite.T(), p.HasNext())
	assert.False(suite.T(), p.HasPrev())
	assert.True(suite.T(), p.HasPages())
}

func (suite *GORM2AdapterTestSuite) TestLastPage() {
	q := suite.db.Model(Post2{})
	p := paginator.New(adapter.NewGORM2Adapter(q), 10)

	p.SetPage(10)
	assert.False(suite.T(), p.HasNext())
	assert.True(suite.T(), p.HasPrev())
}

func (suite *GORM2AdapterTestSuite) TestOutOfRangeCurrentPage() {
	q := suite.db.Model(Post2{})
	p := paginator.New(adapter.NewGORM2Adapter(q), 10)

	var posts []Post2
	p.SetPage(11)
	err := p.Results(&posts)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 10, p.Page())

	posts = make([]Post2, 0)
	p.SetPage(-4)
	assert.True(suite.T(), p.HasNext())
	assert.False(suite.T(), p.HasPrev())
	assert.True(suite.T(), p.HasPages())

	err = p.Results(&posts)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 10)
}

func (suite *GORM2AdapterTestSuite) TestCurrentPageResults() {
	q := suite.db.Model(Post2{})
	p := paginator.New(adapter.NewGORM2Adapter(q), 10)

	var posts []Post2
	p.SetPage(6)
	err := p.Results(&posts)
	assert.NoError(suite.T(), err)

	assert.Len(suite.T(), posts, 10)
	for i, post := range posts {
		assert.Equal(suite.T(), (p.Page()-1)*10+i+1, post.Number)
	}
}

func TestGORM2AdapterTestSuite(t *testing.T) {
	suite.Run(t, new(GORM2AdapterTestSuite))
}
