package lastfm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFormatArgs(t *testing.T) {
	testCases := []struct {
		name string
		in1  P
		in2  P
		out  map[string]string
	}{
		{
			"TestPlainString",
			P{"artist": "a0"},
			P{"plain": []string{"artist"}},
			map[string]string{"artist": "a0"},
		},
		{
			"TestPlainInt",
			P{"id": 29},
			P{"plain": []string{"id"}},
			map[string]string{"id": "29"},
		},
		{
			"TestPlainInt64",
			P{"id": int64(29)},
			P{"plain": []string{"id"}},
			map[string]string{"id": "29"},
		},
		{
			"TestPlainArrayLimit",
			P{"tags": []string{"t0", "t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9", "t10", "t11", "t12"}},
			P{"plain": []string{"tags"}},
			map[string]string{"tags": "t0,t1,t2,t3,t4,t5,t6,t7,t8,t9"},
		},
		{
			"TestPlainArrayOneItem",
			P{"tags": "t"},
			P{"plain": []string{"tags"}},
			map[string]string{"tags": "t"},
		},
		{
			"TestIndexingArray",
			P{"artist": []string{"a0", "a1", "a2"}},
			P{"indexing": []string{"artist"}},
			map[string]string{"artist[0]": "a0", "artist[1]": "a1", "artist[2]": "a2"},
		},
		{
			"TestIndexingArrayOneItem",
			P{"artist": "a"},
			P{"indexing": []string{"artist"}},
			map[string]string{"artist[0]": "a"},
		},
		{
			"TestPlainAndIndexing1",
			P{"artist": []string{"a0", "a1", "a2"}, "recipient": []string{"r0", "r1", "r2"}},
			P{"indexing": []string{"artist"}, "plain": []string{"recipient"}},
			map[string]string{"artist[0]": "a0", "artist[1]": "a1", "artist[2]": "a2", "recipient": "r0,r1,r2"},
		},
		{
			"TestPlainAndIndexing2",
			P{"tags": []string{"t0", "t1", "t2"}, "artist": []string{"a0", "a1"}, "id": 10, "name": "kumakichi"},
			P{"plain": []string{"id", "name", "tags"}, "indexing": []string{"artist"}},
			map[string]string{"tags": "t0,t1,t2", "artist[0]": "a0", "artist[1]": "a1", "id": "10", "name": "kumakichi"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if result, err := formatArgs(tt.in1, tt.in2); !reflect.DeepEqual(tt.out, result) {
				fmt.Printf("result: %+v, error: %v\n", result, err)
				t.Fail()
			}
		})
	}
}
