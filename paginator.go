package paginator

import (
	"errors"
	"math"
)

// DefaultMaxPerPage default number of records per page
const DefaultMaxPerPage = 10

// ErrOutOfRangeCurrentPage current page is out of range error
var ErrOutOfRangeCurrentPage = errors.New("out of range current page")

// Adapter any adapter must implement this interface
type Adapter interface {
	Nums() int
	Slice(offset, length int, data interface{}) error
}

// Paginator structure
type Paginator struct {
	adapter     Adapter
	maxPerPage  int
	CurrentPage int
	nums        int
}

// New paginator constructor
func New(adapter Adapter, maxPerPage int) Paginator {
	if maxPerPage <= 0 {
		maxPerPage = DefaultMaxPerPage
	}

	return Paginator{
		adapter:     adapter,
		maxPerPage:  maxPerPage,
		CurrentPage: 1,
		nums:        -1,
	}
}

// Results stores the current page results into data argument which must be a pointer to a slice.
func (p Paginator) Results(data interface{}) error {
	page := p.CurrentPage
	if page <= 0 {
		page = 1
	}

	return p.pageResults(page, data)
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
	return p.CurrentPage < p.PageNums()
}

// HasPrev returns true if current page is not the first page
func (p Paginator) HasPrev() bool {
	return p.CurrentPage > 1
}

// PageNums returns the total number of pages
func (p Paginator) PageNums() int {
	n := int(math.Ceil(float64(p.Nums()) / float64(p.maxPerPage)))
	if n == 0 {
		n = 1
	}

	return n
}

func (p Paginator) validatePage(page int) error {
	if page <= 0 || page > p.PageNums() {
		return ErrOutOfRangeCurrentPage
	}

	return nil
}

func (p Paginator) pageResults(page int, data interface{}) error {
	err := p.validatePage(page)
	if err != nil {
		return err
	}

	var offset int
	if page > 1 {
		offset = (page - 1) * p.maxPerPage
	}

	return p.adapter.Slice(offset, p.maxPerPage, data)
}
