package main

import (
	"reflect"
	"testing"
)

func TestLengthOfLongestSubstring(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"applebees", 4},
		{"bbbb", 1},
		{"", 0},
		{"z", 1},
		{"makeamericablinkagain", 8},
		{"aabababcabcabcdabcdefff", 6},
	}

	for _, testCase := range testCases {
		if result := lengthOfLongestSubstring(testCase.input); !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected result was incorrect, returned: %v, wanted: %v.", result, testCase.expected)
		}
	}
}
