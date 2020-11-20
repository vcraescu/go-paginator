package view_test

import (
	"github.com/stretchr/testify/suite"
	"github.com/harrifeng/go-paginator"
	"github.com/harrifeng/go-paginator/adapter"
	"github.com/harrifeng/go-paginator/view"
	"testing"
)

type DefaultViewTestSuite struct {
	suite.Suite

	paginator paginator.Paginator
	view      view.Viewer
}

func (suite *DefaultViewTestSuite) SetupTest() {
	data := make([]int, 150)
	for i := 1; i <= 150; i++ {
		data[i-1] = i
	}

	p := paginator.New(adapter.NewSliceAdapter(data), 10)
	suite.paginator = p

	suite.view = view.New(suite.paginator)
}

func (suite *DefaultViewTestSuite) TestItems1stCurrentPage() {
	require := suite.Require()
	pages, err := suite.view.Pages()
	require.NoError(err)
	require.Len(pages, 10)

	_, err = suite.view.Prev()
	require.Equal(paginator.ErrNoPrevPage, err)

	next, err := suite.view.Next()
	require.NoError(err)
	require.Equal(2, next)

	for i, page := range pages {
		require.Equal(i+1, page)
	}
}

func (suite *DefaultViewTestSuite) TestItems3rdCurrentPage() {
	require := suite.Require()
	suite.paginator.SetPage(3)
	pages, err := suite.view.Pages()
	require.NoError(err)

	require.Len(pages, 10)

	prev, err := suite.view.Prev()
	require.NoError(err)
	require.Equal(2, prev)

	next, err := suite.view.Next()
	require.NoError(err)
	require.Equal(4, next)

	for i, page := range pages {
		require.Equal(i+1, page)
	}
}

func (suite *DefaultViewTestSuite) TestItems6thCurrentPage() {
	require := suite.Require()
	suite.paginator.SetPage(6)

	pages, err := suite.view.Pages()
	require.NoError(err)
	require.Len(pages, 10)

	prev, err := suite.view.Prev()
	require.NoError(err)
	require.Equal(5, prev)

	next, err := suite.view.Next()
	require.NoError(err)
	require.Equal(7, next)

	for i, page := range pages {
		require.Equal(i+1, page)
	}
}

func (suite *DefaultViewTestSuite) TestItems7thCurrentPage() {
	require := suite.Require()
	suite.paginator.SetPage(7)

	pages, err := suite.view.Pages()
	require.NoError(err)
	require.Len(pages, 10)

	prev, err := suite.view.Prev()
	require.NoError(err)
	require.Equal(6, prev)

	next, err := suite.view.Next()
	require.NoError(err)
	require.Equal(8, next)

	for i, page := range pages {
		require.Equal(i+2, page)
	}
}

func (suite *DefaultViewTestSuite) TestItems10thCurrentPage() {
	require := suite.Require()
	suite.paginator.SetPage(10)

	pages, err := suite.view.Pages()
	require.NoError(err)
	require.Len(pages, 10)

	prev, err := suite.view.Prev()
	require.NoError(err)
	require.Equal(9, prev)

	next, err := suite.view.Next()
	require.NoError(err)
	require.Equal(11, next)

	for i, page := range pages {
		require.Equal(i+5, page)
	}
}

func (suite *DefaultViewTestSuite) TestItemsLastCurrentPage() {
	require := suite.Require()
	suite.paginator.SetPage(15)

	pages, err := suite.view.Pages()
	require.NoError(err)
	require.Len(pages, 10)

	prev, err := suite.view.Prev()
	require.Equal(14, prev)

	_, err = suite.view.Next()
	require.Equal(paginator.ErrNoNextPage, err)

	for i, page := range pages {
		require.Equal(i+10, page)
	}
}

func (suite *DefaultViewTestSuite) TestItemsNumberOfPagesLessThanLength() {
	data := make([]int, 70)
	for i := 1; i <= 70; i++ {
		data[i-1] = i
	}

	p := paginator.New(adapter.NewSliceAdapter(data), 10)
	suite.paginator = p
	suite.view = view.New(suite.paginator)
	suite.paginator.SetPage(4)

	require := suite.Require()
	pages, err := suite.view.Pages()
	require.NoError(err)
	require.Len(pages, 7)

	prev, err := suite.view.Prev()
	require.NoError(err)
	require.Equal(3, prev)

	next, err := suite.view.Next()
	require.NoError(err)
	require.Equal(5, next)

	for i, page := range pages {
		require.Equal(i+1, page)
	}
}

func (suite *DefaultViewTestSuite) TestItemsInvalidCurrentPage() {
	require := suite.Require()
	suite.paginator.SetPage(-14)

	pages, err := suite.view.Pages()
	require.NoError(err)
	require.Len(pages, 10)

	_, err = suite.view.Prev()
	require.Equal(paginator.ErrNoPrevPage, err)

	next, err := suite.view.Next()
	require.NoError(err)
	require.Equal(2, next)

	for i, page := range pages {
		require.Equal(i+1, page)
	}
}

func TestDefaultViewTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultViewTestSuite))
}
