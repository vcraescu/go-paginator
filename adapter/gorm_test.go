package adapter_test

import (
	"github.com/stretchr/testify/suite"
	"github.com/harrifeng/go-paginator"
	"github.com/harrifeng/go-paginator/adapter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type (
	Post struct {
		ID     uint `gorm:"primary_key"`
		Number int
	}

	GORMAdapterTestSuite struct {
		suite.Suite
		db *gorm.DB
	}
)

func (suite *GORMAdapterTestSuite) SetupTest() {
	require := suite.Require()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(err)

	suite.db = db
	require.NoError(suite.db.AutoMigrate(&Post{}))

	for i := 1; i <= 100; i++ {
		p := Post{
			Number: i,
		}

		require.NoError(suite.db.Save(&p).Error)
	}
}

func (suite *GORMAdapterTestSuite) TearDownTest() {
	require := suite.Require()
	rawDB, err := suite.db.DB()

	require.NoError(err)
	require.NoError(rawDB.Close())
}

func (suite *GORMAdapterTestSuite) TestFirstPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	require := suite.Require()
	pn, err := p.PageNums()
	require.NoError(err)
	require.Equal(10, pn)

	page, err := p.Page()
	require.NoError(err)
	require.Equal(1, page)

	hn, err := p.HasNext()
	require.NoError(err)
	require.True(hn)

	hp, err := p.HasPrev()
	require.NoError(err)
	require.False(hp)

	hpages, err := p.HasPages()
	require.NoError(err)
	require.True(hpages)
}

func (suite *GORMAdapterTestSuite) TestLastPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	p.SetPage(10)

	require := suite.Require()
	hn, err := p.HasNext()
	require.NoError(err)
	require.False(hn)

	hp, err := p.HasPrev()
	require.NoError(err)
	require.True(hp)
}

func (suite *GORMAdapterTestSuite) TestOutOfRangeCurrentPage() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	var posts []Post
	p.SetPage(11)

	require := suite.Require()

	err := p.Results(&posts)
	require.NoError(err)

	page, err := p.Page()
	require.NoError(err)
	require.Equal(10, page)

	posts = make([]Post, 0)
	p.SetPage(-4)

	hn, err := p.HasNext()
	require.NoError(err)
	require.True(hn)

	hp, err := p.HasPrev()
	require.NoError(err)
	require.False(hp)

	hpages, err := p.HasPages()
	require.NoError(err)
	require.True(hpages)

	err = p.Results(&posts)
	require.NoError(err)
	require.Len(posts, 10)
}

func (suite *GORMAdapterTestSuite) TestCurrentPageResults() {
	q := suite.db.Model(Post{})
	p := paginator.New(adapter.NewGORMAdapter(q), 10)

	var posts []Post
	p.SetPage(6)

	require := suite.Require()
	require.NoError(p.Results(&posts))

	require.Len(posts, 10)
	for i, post := range posts {
		page, err := p.Page()
		require.NoError(err)
		require.Equal((page-1)*10+i+1, post.Number)
	}
}

func TestGORMAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(GORMAdapterTestSuite))
}
