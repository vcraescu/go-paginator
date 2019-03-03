package adapter_test

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"testing"
)

type ArrayAdapterTestSuite struct {
	suite.Suite
	data []int
}

func (suite *ArrayAdapterTestSuite) SetupTest() {
	suite.data = make([]int, 100)
	for i := 1; i <= 100; i++ {
		suite.data[i-1] = i
	}
}

func (suite *ArrayAdapterTestSuite) TearDownTest() {
	suite.data = make([]int, 0)
}

func (suite *ArrayAdapterTestSuite) TestFirstPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	assert.Equal(suite.T(), 10, p.PageNums())
	assert.Equal(suite.T(), 1, p.Page())
	assert.True(suite.T(), p.HasNext())
	assert.False(suite.T(), p.HasPrev())
	assert.True(suite.T(), p.HasPages())
}

func (suite *ArrayAdapterTestSuite) TestLastPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	p.SetPage(10)
	assert.False(suite.T(), p.HasNext())
	assert.True(suite.T(), p.HasPrev())
}

func (suite *ArrayAdapterTestSuite) TestOutOfRangeCurrentPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	var pages []int
	p.SetPage(11)
	err := p.Results(&pages)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 10, p.Page())

	pages = make([]int, 0)
	p.SetPage(-4)
	assert.True(suite.T(), p.HasNext())
	assert.False(suite.T(), p.HasPrev())
	assert.True(suite.T(), p.HasPages())

	err = p.Results(&pages)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), pages, 10)
}

func (suite *ArrayAdapterTestSuite) TestCurrentPageResults() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	var pages []int
	p.SetPage(6)
	err := p.Results(&pages)
	assert.NoError(suite.T(), err)

	assert.Len(suite.T(), pages, 10)
	for i, page := range pages {
		assert.Equal(suite.T(), (p.Page()-1)*10+i+1, page)
	}
}

func TestArrayAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(ArrayAdapterTestSuite))
}
