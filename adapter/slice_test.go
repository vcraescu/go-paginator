package adapter_test

import (
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

func (suite *ArrayAdapterTestSuite) TestFirstPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

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

func (suite *ArrayAdapterTestSuite) TestLastPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	p.SetPage(10)

	require := suite.Require()

	hn, err := p.HasNext()
	require.NoError(err)
	require.False(hn)

	hp, err := p.HasPrev()
	require.NoError(err)
	require.True(hp)
}

func (suite *ArrayAdapterTestSuite) TestOutOfRangeCurrentPage() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	require := suite.Require()

	var pages []int
	p.SetPage(11)
	err := p.Results(&pages)
	require.NoError(err)

	page, err := p.Page()
	require.NoError(err)
	require.Equal(10, page)

	pages = make([]int, 0)
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

	err = p.Results(&pages)
	require.NoError(err)
	require.Len(pages, 10)
}

func (suite *ArrayAdapterTestSuite) TestCurrentPageResults() {
	p := paginator.New(adapter.NewSliceAdapter(suite.data), 10)

	require := suite.Require()
	var pages []int
	p.SetPage(6)

	require.NoError(p.Results(&pages))

	require.Len(pages, 10)
	for i, page := range pages {
		expectedPage, err := p.Page()
		require.NoError(err)
		require.Equal((expectedPage-1)*10+i+1, page)
	}
}

func TestArrayAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(ArrayAdapterTestSuite))
}
