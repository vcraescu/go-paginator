package paginator_test

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vcraescu/go-paginator"
	"testing"
)

type Post struct {
	ID     uint `gorm:"primary_key"`
	Number int
}

type GenericAdapter struct {
	nums int
}

func (p GenericAdapter) Nums() int {
	return p.nums
}

func (p GenericAdapter) Slice(offset, length int, data interface{}) error {
	s := data.(*[]Post)

	for n := offset + 1; n < offset+length+1; n++ {
		*s = append(*s, Post{Number: n})
	}

	return nil
}

type PaginatorTestSuite struct {
	suite.Suite
}

func (suite *PaginatorTestSuite) TestFirstPage() {
	p := paginator.New(&GenericAdapter{nums: 100}, 10)

	assert.Equal(suite.T(), 10, p.PageNums())
	assert.Equal(suite.T(), 1, p.Page())
	assert.True(suite.T(), p.HasNext())
	assert.False(suite.T(), p.HasPrev())
	assert.True(suite.T(), p.HasPages())
}

func (suite *PaginatorTestSuite) TestLastPage() {
	p := paginator.New(&GenericAdapter{nums: 100}, 10)

	p.SetPage(10)
	assert.False(suite.T(), p.HasNext())
	assert.True(suite.T(), p.HasPrev())
}

func (suite *PaginatorTestSuite) TestOutOfRangeCurrentPage() {
	p := paginator.New(&GenericAdapter{nums: 100}, 10)

	var posts []Post
	p.SetPage(11)
	err := p.Results(&posts)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 10, p.Page())

	posts = make([]Post, 0)
	p.SetPage(-4)
	assert.Equal(suite.T(), 1, p.Page())
	assert.True(suite.T(), p.HasNext())
	assert.False(suite.T(), p.HasPrev())
	assert.True(suite.T(), p.HasPages())

	err = p.Results(&posts)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 10)
}

func (suite *PaginatorTestSuite) TestCurrentPageResults() {
	p := paginator.New(&GenericAdapter{nums: 100}, 10)

	var posts []Post
	p.SetPage(6)
	err := p.Results(&posts)
	assert.NoError(suite.T(), err)

	assert.Len(suite.T(), posts, 10)
	for i, post := range posts {
		assert.Equal(suite.T(), (p.Page()-1)*10+i+1, post.Number)
	}
}

func TestPluginTestSuite(t *testing.T) {
	suite.Run(t, new(PaginatorTestSuite))
}
