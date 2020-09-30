package adapter

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestIsSlice(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "slice",
			value:    []string{},
			expected: true,
		},
		{
			name:  "nil",
			value: nil,
		},
		{
			name:  "non-slice",
			value: 4,
		},
		{
			name:     "ptr slice",
			value:    &[]string{},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, isSlice(test.value))
		})
	}
}

func TestIsPtr(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "pointer",
			value:    &[]string{},
			expected: true,
		},
		{
			name:  "non pointer",
			value: 4,
		},
		{
			name:  "nil",
			value: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, isSlice(test.value))
		})
	}
}

func TestMakeSlice(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		length int
		err    bool
	}{
		{
			name:   "slice",
			value:  &[]string{},
			length: 10,
		},
		{
			name:  "non slice",
			value: 4,
			err:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := makeSlice(test.value, test.length, test.length)
			if test.err {
				require.Error(t, err)
				return
			}

			require.Equal(t, test.length, reflect.Indirect(reflect.ValueOf(test.value)).Len())
		})
	}
}
