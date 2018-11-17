package main

import (
	"reflect"
	"testing"
)

func TestFindAddend(t *testing.T) {
	tests := []struct {
		targetSum int
		input     []int
		expected  []int
	}{
		{
			targetSum: 9,
			input:     []int{2, 7, 11, 15},
			expected:  []int{0, 1},
		},
		{
			targetSum: 9,
			input:     []int{0, -40, 1, -1, 19, -3, 5, 90, -2, 7, 11, 15},
			expected:  []int{8, 10},
		},
		{
			targetSum: -4,
			input:     []int{0, -40, 1, -1, 19, -2, 5, 90, -2, 7, 11, 15},
			expected:  []int{5, 8},
		},
		{
			targetSum: 0,
			input:     []int{0, 2, 7, 11, 15, 0},
			expected:  []int{0, 5},
		},
		{
			targetSum: 0,
			input:     []int{0, -12, 7, 11, 15, -3, 5, 12},
			expected:  []int{1, 7},
		},
		{
			targetSum: -10,
			input:     []int{4, 0, -12, 2, 7, 11, 15, -3, 5, 12},
			expected:  []int{2, 3},
		},
	}

	for _, test := range tests {
		if result := findAddend(test.input, test.targetSum); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected result was incorrect, returned: %v, wanted: %v.", result, test.expected)
		}
	}
}
