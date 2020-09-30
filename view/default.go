package view

import (
	"github.com/vcraescu/go-paginator"
)

type (
	// Viewer interface
	Viewer interface {
		Pages() ([]int, error)
		Next() (int, error)
		Prev() (int, error)
		Last() (int, error)
		Current() (int, error)
	}

	// DefaultView viewer interface implementation
	// The paginator will look like the one from google
	DefaultView struct {
		Paginator paginator.Paginator
		Proximity int
	}
)

// New DefaultView constructor
func New(p paginator.Paginator) Viewer {
	return &DefaultView{
		Paginator: p,
		Proximity: 5,
	}
}

// Next returns next page number or zero if current page is the last page
func (v *DefaultView) Next() (int, error) {
	return v.Paginator.NextPage()
}

// Prev returns previous page number or zero if current page is first page
func (v *DefaultView) Prev() (int, error) {
	return v.Paginator.PrevPage()
}

// Last returns last page number
func (v *DefaultView) Last() (int, error) {
	return v.Paginator.PageNums()
}

// Current returns current page number
func (v *DefaultView) Current() (int, error) {
	return v.Paginator.Page()
}

// Pages returns the list of pages
func (v *DefaultView) Pages() ([]int, error) {
	var items []int
	hasPages, err := v.Paginator.HasPages()
	if err != nil {
		return nil, err
	}

	if !hasPages {
		return items, nil
	}

	items = make([]int, 0)
	length := v.Proximity * 2
	pn, err := v.Paginator.PageNums()
	if err != nil {
		return nil, err
	}

	if pn < length {
		pn, err := v.Paginator.PageNums()
		if err != nil {
			return nil, err
		}

		length = pn
	}

	proximityLeft := length / 2
	proximityRight := (length / 2) - 1
	if length%2 != 0 {
		proximityRight = proximityLeft
	}

	page, err := v.Paginator.Page()
	if err != nil {
		return nil, err
	}

	start := page - proximityLeft
	end := page + proximityRight
	if start <= 0 {
		start = 1
		end = length
	}

	for page = start; page <= end; page++ {
		items = append(items, page)
	}

	return items, nil
}
