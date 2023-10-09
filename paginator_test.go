package paginator_test

import (
	"testing"

	"github.com/harrifeng/go-paginator"
	"github.com/stretchr/testify/suite"
)

var (
	_ paginator.Adapter = (*GenericAdapter)(nil)
)

type (
	Post struct {
		ID     uint `gorm:"primary_key"`
		Number int
	}

	GenericAdapter struct {
		nums int64
	}
)

func (p GenericAdapter) Nums() (int64, error) {
	return p.nums, nil
}

func (p GenericAdapter) Slice(order string, offset, length int, data interface{}) error {
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

	pn, _ := p.PageNums()
	suite.Equal(10, pn)

	page, _ := p.Page()
	suite.Equal(1, page)

	hn, _ := p.HasNext()
	suite.True(hn)

	hp, _ := p.HasPrev()
	suite.False(hp)

	hpages, _ := p.HasPages()
	suite.True(hpages)
}

func (suite *PaginatorTestSuite) TestLastPage() {
	p := paginator.New(&GenericAdapter{nums: 100}, 10)

	p.SetPage(10)

	hn, _ := p.HasNext()
	suite.False(hn)

	hp, _ := p.HasPrev()
	suite.True(hp)
}

func (suite *PaginatorTestSuite) TestOutOfRangeCurrentPage() {
	p := paginator.New(&GenericAdapter{nums: 100}, 10)

	var posts []Post
	p.SetPage(11)
	err := p.Results(&posts)
	suite.NoError(err)

	page, _ := p.Page()
	suite.Equal(10, page)

	posts = make([]Post, 0)
	p.SetPage(-4)

	page, _ = p.Page()
	suite.Equal(1, page)

	hn, _ := p.HasNext()
	suite.True(hn)

	hp, _ := p.HasPrev()
	suite.False(hp)

	hpages, _ := p.HasPages()
	suite.True(hpages)

	err = p.Results(&posts)
	suite.NoError(err)
	suite.Len(posts, 10)
}

func (suite *PaginatorTestSuite) TestCurrentPageResults() {
	p := paginator.New(&GenericAdapter{nums: 100}, 10)

	var posts []Post
	p.SetPage(6)
	err := p.Results(&posts)
	suite.NoError(err)

	suite.Len(posts, 10)
	for i, post := range posts {
		page, _ := p.Page()
		suite.Equal((page-1)*10+i+1, post.Number)
	}
}

func TestPluginTestSuite(t *testing.T) {
	suite.Run(t, new(PaginatorTestSuite))
}
