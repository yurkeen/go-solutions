package main

import (
	"reflect"
	"testing"
)

func TestLengthOfLongestSubstring(t *testing.T) {
	tests := []struct {
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

	for _, test := range tests {
		if result := lengthOfLongestSubstring(test.input); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected result was incorrect, returned: %v, wanted: %v.", result, test.expected)
		}
	}
}
