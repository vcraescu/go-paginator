package paginator

import (
	"errors"
	"fmt"
	"math"
)

// DefaultMaxPerPage default number of records per page
const DefaultMaxPerPage = 10

var (
	// ErrNoPrevPage current page is first page
	ErrNoPrevPage = errors.New("no previous page")

	// ErrNoNextPage current page is last page
	ErrNoNextPage = errors.New("no next page")
)

type (
	// Adapter any adapter must implement this interface
	Adapter interface {
		Nums() (int64, error)
		Slice(order string, offset, length int, data interface{}) error
	}

	// Paginator interface
	Paginator interface {
		SetPage(page int)
		Page() (int, error)
		Results(data interface{}) error
		Nums() (int64, error)
		HasPages() (bool, error)
		HasNext() (bool, error)
		PrevPage() (int, error)
		NextPage() (int, error)
		HasPrev() (bool, error)
		PageNums() (int, error)
		PerPage() (int, error)
		SetSort(sort string)
		Sort() (string, error)
	}

	// Paginator structure
	paginator struct {
		adapter    Adapter
		maxPerPage int
		page       int
		nums       int64
		sort       string
	}
)

// New paginator constructor
func New(adapter Adapter, maxPerPage int) Paginator {
	if maxPerPage <= 0 {
		maxPerPage = DefaultMaxPerPage
	}

	return &paginator{
		adapter:    adapter,
		maxPerPage: maxPerPage,
		page:       1,
		nums:       -1,
		sort:       "+id",
	}
}

func (p *paginator) SetSort(sort string) {

	p.sort = sort
}

func (p *paginator) Sort() (string, error) {
	if len(p.sort) < 2 {
		return "", fmt.Errorf("too short sort parameter")
	}

	if p.sort[0] == '+' {
		return fmt.Sprintf("%v asc", p.sort[1:len(p.sort)]), nil
	} else if p.sort[0] == '-' {
		return fmt.Sprintf("%v desc", p.sort[1:len(p.sort)]), nil
	} else {
		return "", fmt.Errorf("sort parameter should start with `-` as desc or `+` as asc")
	}

}

// SetPage set current page
func (p *paginator) SetPage(page int) {
	if page <= 0 {
		page = 1
	}

	p.page = page
}

// Page returns current page
func (p paginator) Page() (int, error) {
	pn, err := p.PageNums()
	if err != nil {
		return 0, err
	}

	if p.page > pn {
		return pn, nil
	}

	return p.page, nil
}

// Results stores the current page results into data argument which must be a pointer to a slice.
func (p paginator) Results(data interface{}) error {
	var offset int
	page, err := p.Page()
	if err != nil {
		return err
	}

	if page > 1 {
		offset = (page - 1) * p.maxPerPage
	}

	sortStr, err := p.Sort()
	if err != nil {
		return err
	}

	return p.adapter.Slice(sortStr, offset, p.maxPerPage, data)
}

// Nums returns the total number of records
func (p *paginator) Nums() (int64, error) {
	var err error
	if p.nums == -1 {
		p.nums, err = p.adapter.Nums()
		if err != nil {
			return 0, err
		}
	}

	return p.nums, nil
}

// HasPages returns true if there is more than one page
func (p paginator) HasPages() (bool, error) {
	n, err := p.Nums()
	if err != nil {
		return false, err
	}

	return n > int64(p.maxPerPage), nil
}

// HasNext returns true if current page is not the last page
func (p paginator) HasNext() (bool, error) {
	pn, err := p.PageNums()
	if err != nil {
		return false, err
	}

	page, err := p.Page()
	if err != nil {
		return false, err
	}

	return page < pn, nil
}

// PrevPage returns previous page number or ErrNoPrevPage if current page is first page
func (p paginator) PrevPage() (int, error) {
	hp, err := p.HasPrev()
	if err != nil {
		return 0, nil
	}

	if !hp {
		return 0, ErrNoPrevPage
	}

	page, err := p.Page()
	if err != nil {
		return 0, err
	}

	return page - 1, nil
}

// NextPage returns next page number or ErrNoNextPage if current page is last page
func (p paginator) NextPage() (int, error) {
	hn, err := p.HasNext()
	if err != nil {
		return 0, err
	}

	if !hn {
		return 0, ErrNoNextPage
	}

	page, err := p.Page()
	if err != nil {
		return 0, err
	}

	return page + 1, nil
}

// HasPrev returns true if current page is not the first page
func (p paginator) HasPrev() (bool, error) {
	page, err := p.Page()
	if err != nil {
		return false, err
	}

	return page > 1, nil
}

// PageNums returns the total number of pages
func (p paginator) PageNums() (int, error) {
	n, err := p.Nums()
	if err != nil {
		return 0, err
	}

	n = int64(math.Ceil(float64(n) / float64(p.maxPerPage)))
	if n == 0 {
		n = 1
	}

	return int(n), nil
}

// PerPage returns the exact per page for the pagination.
func (p paginator) PerPage() (int, error) {
	return p.maxPerPage, nil
}
