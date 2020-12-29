package slices_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alexfacciorusso/winget-generate/slices"
)

func TestElementToFirst(t *testing.T) {
	var testCases = []struct {
		slice    []string
		element  string
		expected []string
	}{
		{[]string{"one", "two", "three"}, "two", []string{"two", "one", "three"}},
		{[]string{"one", "two", "three"}, "three", []string{"three", "one", "two"}},
		{[]string{"one", "two", "three"}, "four", []string{"one", "two", "three"}},
		{[]string{"one"}, "one", []string{"one"}},
		{[]string{"one"}, "two", []string{"one"}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			result := slices.ElementToFirst(tc.slice, tc.element)

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Case %+v not satisfied. Result: %v", tc, result)
			}
		})
	}
}
