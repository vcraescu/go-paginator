package adapter

import (
	"fmt"
	"reflect"

	"github.com/harrifeng/go-paginator"
)

// SliceAdapter slice adapter to be passed to paginator constructor to paginate a slice of elements.
type SliceAdapter struct {
	src interface{}
}

// NewSliceAdapter slice adapter construct receive the slice source which needs to be paginated.
func NewSliceAdapter(source interface{}) paginator.Adapter {
	if isPtr(source) || !isSlice(source) {
		panic(fmt.Sprintf("expected slice but got %s", reflect.TypeOf(source).Kind()))
	}

	return &SliceAdapter{source}
}

// Nums returns the number of elements
func (a *SliceAdapter) Nums() (int64, error) {
	n := reflect.ValueOf(a.src).Len()

	return int64(n), nil
}

// Slice stores into dest argument a slice of the results.
// dest argument must be a pointer to a slice
func (a *SliceAdapter) Slice(order string, offset, length int, dest interface{}) error {
	// adjust the length for the last page
	va := reflect.ValueOf(a.src)
	totalsize := va.Len()
	if totalsize < length+offset {
		length = totalsize - offset
	}

	if err := makeSlice(dest, length, length); err != nil {
		return err
	}

	vs := va.Slice(offset, offset+length)
	vt := reflect.ValueOf(dest).Elem()
	for i := 0; i < vs.Len(); i++ {
		vt.Index(i).Set(reflect.ValueOf(vs.Index(i).Interface()))
	}

	return nil
}
