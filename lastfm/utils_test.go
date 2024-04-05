package lastfm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestToString(t *testing.T) {
	testCases := []struct {
		name          string
		input         interface{}
		expectedStr   string
		expectedError error
	}{
		{
			name:          "String",
			input:         "hello",
			expectedStr:   "hello",
			expectedError: nil,
		},
		{
			name:          "Integer",
			input:         42,
			expectedStr:   "42",
			expectedError: nil,
		},
		{
			name:          "Integer 64",
			input:         int64(123),
			expectedStr:   "123",
			expectedError: nil,
		},
		{
			name:          "Slice of strings",
			input:         []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"},
			expectedStr:   "a,b,c,d,e,f,g,h,i,j",
			expectedError: nil,
		},
		{
			name:          "Unsupported Bool",
			input:         true,
			expectedStr:   "",
			expectedError: newLibError(ErrorInvalidTypeOfArgument, Messages[ErrorInvalidTypeOfArgument]),
		},
		{
			name:          "Unsupported nil",
			input:         nil,
			expectedStr:   "",
			expectedError: newLibError(ErrorInvalidTypeOfArgument, Messages[ErrorInvalidTypeOfArgument]),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			str, err := toString(tc.input)

			if str != tc.expectedStr {
				t.Errorf("Expected string to be %s, but got %s", tc.expectedStr, str)
			}

			if tc.expectedError != nil && err == nil && tc.expectedError.Error() != err.Error() {
				t.Errorf("Expected error to be %v, but got %v", tc.expectedError, err)
			}
		})
	}
}

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
