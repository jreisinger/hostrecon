package hostrecon

import (
	"reflect"
	"testing"
)

func TestNelems(t *testing.T) {
	var testcases = []struct {
		ss            []string
		n             int
		wantElems     []string
		wantRemaining int
	}{
		{
			ss:            []string{"apple", "banana", "cherry", "date", "elderberry"},
			n:             3,
			wantElems:     []string{"apple", "banana", "cherry"},
			wantRemaining: 2,
		},
		{
			// panic: runtime error: slice bounds out of range [:1] with capacity 0
			ss:            []string{},
			n:             1,
			wantElems:     []string{},
			wantRemaining: 0,
		},
	}

	for _, tc := range testcases {
		gotElems, gotRemaining := maxElems(tc.ss, tc.n)
		if !reflect.DeepEqual(tc.wantElems, gotElems) {
			t.Errorf("want elems %v, got %v", tc.wantElems, gotElems)
		}
		if tc.wantRemaining != gotRemaining {
			t.Errorf("want remaining %d, got %d", tc.wantRemaining, gotRemaining)
		}
	}
}
