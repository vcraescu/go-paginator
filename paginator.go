package paginator

import (
	"errors"
	"math"
)

// DefaultMaxPerPage default number of records per page
const DefaultMaxPerPage = 10

// ErrNoPrevPage current page is first page
var ErrNoPrevPage = errors.New("no previous page")

// ErrNoNextPage current page is last page
var ErrNoNextPage = errors.New("no next page")

// Adapter any adapter must implement this interface
type Adapter interface {
	Nums() int
	Slice(offset, length int, data interface{}) error
}

// Paginator structure
type Paginator struct {
	adapter    Adapter
	maxPerPage int
	page       int
	nums       int
}

// New paginator constructor
func New(adapter Adapter, maxPerPage int) Paginator {
	if maxPerPage <= 0 {
		maxPerPage = DefaultMaxPerPage
	}

	return Paginator{
		adapter:    adapter,
		maxPerPage: maxPerPage,
		page:       1,
		nums:       -1,
	}
}

// SetPage set current page
func (p *Paginator) SetPage(page int) {
	if page <= 0 {
		page = 1
	}

	p.page = page
}

// Page returns current page
func (p Paginator) Page() int {
	pn := p.PageNums()
	if p.page > pn {
		return pn
	}

	return p.page
}

// Results stores the current page results into data argument which must be a pointer to a slice.
func (p Paginator) Results(data interface{}) error {
	var offset int
	page := p.Page()
	if page > 1 {
		offset = (page - 1) * p.maxPerPage
	}

	return p.adapter.Slice(offset, p.maxPerPage, data)
}

// Nums returns the total number of records
func (p *Paginator) Nums() int {
	if p.nums == -1 {
		p.nums = p.adapter.Nums()
	}

	return p.nums
}

// HasPages returns true if there is more than one page
func (p Paginator) HasPages() bool {
	return p.Nums() > p.maxPerPage
}

// HasNext returns true if current page is not the last page
func (p Paginator) HasNext() bool {
	return p.Page() < p.PageNums()
}

// PrevPage returns previous page number or ErrNoPrevPage if current page is first page
func (p Paginator) PrevPage() (int, error) {
	if !p.HasPrev() {
		return 0, ErrNoPrevPage
	}

	return p.Page() - 1, nil
}

// NextPage returns next page number or ErrNoNextPage if current page is last page
func (p Paginator) NextPage() (int, error) {
	if !p.HasNext() {
		return 0, ErrNoNextPage
	}

	return p.Page() + 1, nil
}

// HasPrev returns true if current page is not the first page
func (p Paginator) HasPrev() bool {
	return p.Page() > 1
}

// PageNums returns the total number of pages
func (p Paginator) PageNums() int {
	n := int(math.Ceil(float64(p.Nums()) / float64(p.maxPerPage)))
	if n == 0 {
		n = 1
	}

	return n
}
