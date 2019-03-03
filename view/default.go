package view

import (
	"github.com/vcraescu/go-paginator"
)

// Viewer interface
type Viewer interface {
	Pages() []int
}

// DefaultView viewer interface implementation
// The paginator will look like the one from google
type DefaultView struct {
	Paginator *paginator.Paginator
	Proximity int
}

// New DefaultView constructor
func New(p *paginator.Paginator) DefaultView {
	return DefaultView{
		Paginator: p,
		Proximity: 5,
	}
}

// Next returns next page number or zero if current page is the last page
func (v DefaultView) Next() int {
	page, _ := v.Paginator.NextPage()

	return page
}

// Prev returns previous page number or zero if current page is first page
func (v DefaultView) Prev() int {
	page, _ := v.Paginator.PrevPage()

	return page
}

// Last returns last page number
func (v DefaultView) Last() int {
	return v.Paginator.PageNums()
}

// Current returns current page number
func (v DefaultView) Current() int {
	return v.Paginator.Page()
}

// Pages returns the list of pages
func (v DefaultView) Pages() []int {
	var items []int
	if !v.Paginator.HasPages() {
		return items
	}

	items = make([]int, 0)
	length := v.Proximity * 2
	if v.Paginator.PageNums() < length {
		length = v.Paginator.PageNums()
	}

	proximityLeft := length / 2
	proximityRight := (length / 2) - 1
	if length%2 != 0 {
		proximityRight = proximityLeft
	}

	start := v.Paginator.Page() - proximityLeft
	end := v.Paginator.Page() + proximityRight
	if start <= 0 {
		start = 1
		end = length
	}

	for page := start; page <= end; page++ {
		items = append(items, page)
	}

	return items
}
