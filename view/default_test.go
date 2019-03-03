package view_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
	"github.com/vcraescu/go-paginator/view"
	"testing"
)

type DefaultViewTestSuite struct {
	suite.Suite
	paginator *paginator.Paginator
	view      view.DefaultView
}

func (suite *DefaultViewTestSuite) SetupTest() {
	data := make([]int, 150)
	for i := 1; i <= 150; i++ {
		data[i-1] = i
	}

	p := paginator.New(adapter.NewSliceAdapter(data), 10)
	suite.paginator = &p

	suite.view = view.New(suite.paginator)
}

func (suite *DefaultViewTestSuite) TestItems1stCurrentPage() {
	pages := suite.view.Pages()

	assert.Len(suite.T(), pages, 10)
	assert.Equal(suite.T(), 0, suite.view.Prev())
	assert.Equal(suite.T(), 2, suite.view.Next())

	for i, page := range pages {
		assert.Equal(suite.T(), i+1, page)
	}
}

func (suite *DefaultViewTestSuite) TestItems3rdCurrentPage() {
	suite.paginator.SetPage(3)
	pages := suite.view.Pages()

	assert.Len(suite.T(), pages, 10)
	assert.Equal(suite.T(), 2, suite.view.Prev())
	assert.Equal(suite.T(), 4, suite.view.Next())

	for i, page := range pages {
		assert.Equal(suite.T(), i+1, page)
	}
}

func (suite *DefaultViewTestSuite) TestItems6thCurrentPage() {
	suite.paginator.SetPage(6)
	pages := suite.view.Pages()

	assert.Len(suite.T(), pages, 10)
	assert.Equal(suite.T(), 5, suite.view.Prev())
	assert.Equal(suite.T(), 7, suite.view.Next())

	for i, page := range pages {
		assert.Equal(suite.T(), i+1, page)
	}
}

func (suite *DefaultViewTestSuite) TestItems7thCurrentPage() {
	suite.paginator.SetPage(7)
	pages := suite.view.Pages()

	assert.Len(suite.T(), pages, 10)
	assert.Equal(suite.T(), 6, suite.view.Prev())
	assert.Equal(suite.T(), 8, suite.view.Next())

	for i, page := range pages {
		assert.Equal(suite.T(), i+2, page)
	}
}

func (suite *DefaultViewTestSuite) TestItems10thCurrentPage() {
	suite.paginator.SetPage(10)
	pages := suite.view.Pages()

	assert.Len(suite.T(), pages, 10)
	assert.Equal(suite.T(), 9, suite.view.Prev())
	assert.Equal(suite.T(), 11, suite.view.Next())

	for i, page := range pages {
		assert.Equal(suite.T(), i+5, page)
	}
}

func (suite *DefaultViewTestSuite) TestItemsLastCurrentPage() {
	suite.paginator.SetPage(15)
	pages := suite.view.Pages()

	assert.Len(suite.T(), pages, 10)
	assert.Equal(suite.T(), 14, suite.view.Prev())
	assert.Equal(suite.T(), 0, suite.view.Next())

	for i, page := range pages {
		assert.Equal(suite.T(), i+10, page)
	}
}

func (suite *DefaultViewTestSuite) TestItemsNumberOfPagesLessThanLength() {
	data := make([]int, 70)
	for i := 1; i <= 70; i++ {
		data[i-1] = i
	}

	p := paginator.New(adapter.NewSliceAdapter(data), 10)
	suite.paginator = &p

	suite.view = view.New(suite.paginator)

	suite.paginator.SetPage(4)
	pages := suite.view.Pages()

	assert.Len(suite.T(), pages, 7)
	assert.Equal(suite.T(), 3, suite.view.Prev())
	assert.Equal(suite.T(), 5, suite.view.Next())

	for i, page := range pages {
		assert.Equal(suite.T(), i+1, page)
	}
}

func (suite *DefaultViewTestSuite) TestItemsInvalidCurrentPage() {
	suite.paginator.SetPage(-14)
	pages := suite.view.Pages()

	assert.Len(suite.T(), pages, 10)
	assert.Equal(suite.T(), 0, suite.view.Prev())
	assert.Equal(suite.T(), 2, suite.view.Next())

	for i, page := range pages {
		assert.Equal(suite.T(), i+1, page)
	}
}

func TestDefaultViewTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultViewTestSuite))
}
