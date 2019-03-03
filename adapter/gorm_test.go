package adapter_test

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"math/rand"
	"os"
	"testing"
)

var dbName = fmt.Sprintf("test_%d.db", rand.Int())

type Post struct {
	ID     uint `gorm:"primary_key"`
	Number int
}

type GORMAdapterTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *GORMAdapterTestSuite) SetupTest() {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		panic(fmt.Errorf("setup test: %s", err))
	}

	suite.db = db
	suite.db.AutoMigrate(&Post{})

	for i := 1; i <= 100; i++ {
		p := Post{
			Number: i,
		}

		suite.db.Save(&p)
	}
}

func (suite *GORMAdapterTestSuite) TearDownTest() {
	if err := suite.db.Close(); err != nil {
		panic(fmt.Errorf("tear down test: %s", err))
	}

	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("tear down test: %s", err))
	}
}

func (suite *GORMAdapterTestSuite) TestFirstPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	assert.Equal(suite.T(), 10, p.PageNums())
	assert.Equal(suite.T(), 1, p.Page())
	assert.True(suite.T(), p.HasNext())
	assert.False(suite.T(), p.HasPrev())
	assert.True(suite.T(), p.HasPages())
}

func (suite *GORMAdapterTestSuite) TestLastPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	p.SetPage(10)
	assert.False(suite.T(), p.HasNext())
	assert.True(suite.T(), p.HasPrev())
}

func (suite *GORMAdapterTestSuite) TestOutOfRangeCurrentPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	var posts []Post
	p.SetPage(11)
	err := p.Results(&posts)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 10, p.Page())

	posts = make([]Post, 0)
	p.SetPage(-4)
	assert.True(suite.T(), p.HasNext())
	assert.False(suite.T(), p.HasPrev())
	assert.True(suite.T(), p.HasPages())

	err = p.Results(&posts)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 10)
}

func (suite *GORMAdapterTestSuite) TestCurrentPageResults() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	var posts []Post
	p.SetPage(6)
	err := p.Results(&posts)
	assert.NoError(suite.T(), err)

	assert.Len(suite.T(), posts, 10)
	for i, post := range posts {
		assert.Equal(suite.T(), (p.Page()-1)*10+i+1, post.Number)
	}
}

func TestGORMAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(GORMAdapterTestSuite))
}
