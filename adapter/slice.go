package adapter

import (
	"fmt"
	"reflect"
)

// SliceAdapter slice adapter to be passed to paginator constructor to paginate a slice of elements.
type SliceAdapter struct {
	src interface{}
}

// NewSliceAdapter slice adapter construct receive the slice source which needs to be paginated.
func NewSliceAdapter(source interface{}) SliceAdapter {
	if isPtr(source) || !isSlice(source) {
		panic(fmt.Sprintf("expected slice but got %s", reflect.TypeOf(source).Kind()))
	}

	return SliceAdapter{source}
}

// Nums returns the number of elements
func (a SliceAdapter) Nums() int {
	return reflect.ValueOf(a.src).Len()
}

// Slice stores into dest argument a slice of the results.
// dest argument must be a pointer to a slice
func (a SliceAdapter) Slice(offset, length int, dest interface{}) error {
	if !isPtr(dest) {
		panic(fmt.Sprintf("expected slice pointer but got %s", reflect.TypeOf(dest).Kind()))
	}

	indDest := reflect.Indirect(reflect.ValueOf(dest))
	if !isSlice(indDest.Interface()) {
		panic(fmt.Sprintf("expected slice but got %s", indDest.Kind()))
	}
	// adjust the length for the last page
	va := reflect.ValueOf(a.src)
	{
		totalsize := va.Len()
		if totalsize < length+offset {
			length = totalsize - offset
		}
	}
	makeSlice(dest, length, length)
	vs := va.Slice(offset, offset+length)
	vt := reflect.ValueOf(dest).Elem()
	for i := 0; i < vs.Len(); i++ {
		vt.Index(i).Set(reflect.ValueOf(vs.Index(i).Interface()))
	}

	return nil
}

func isSlice(data interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(data))

	return v.Kind() == reflect.Slice
}

func isPtr(data interface{}) bool {
	t := reflect.TypeOf(data)

	return t.Kind() == reflect.Ptr
}

func makeSlice(data interface{}, length, cap int) {
	ind := reflect.Indirect(reflect.ValueOf(data))

	typ := reflect.TypeOf(ind.Interface())
	reflect.ValueOf(data).Elem().Set(reflect.MakeSlice(typ, length, cap))
}
