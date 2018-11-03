package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAddTwoNumbers(t *testing.T) {
	testCases := []struct {
		leftList     *ListNode
		righList     *ListNode
		expectedList *ListNode
	}{
		{
			leftList:     &ListNode{4, &ListNode{5, &ListNode{6, nil}}},
			righList:     &ListNode{3, &ListNode{2, &ListNode{1, nil}}},
			expectedList: &ListNode{7, &ListNode{7, &ListNode{7, nil}}},
		},
		{
			leftList:     &ListNode{2, &ListNode{1, nil}},
			righList:     &ListNode{4, &ListNode{3, &ListNode{2, nil}}},
			expectedList: &ListNode{6, &ListNode{4, &ListNode{2, nil}}},
		},
		{
			leftList:     &ListNode{0, nil},
			righList:     &ListNode{0, &ListNode{1, nil}},
			expectedList: &ListNode{0, &ListNode{1, nil}},
		},
		{
			leftList:     &ListNode{1, &ListNode{0, &ListNode{1, nil}}},
			righList:     &ListNode{9, &ListNode{9, nil}},
			expectedList: &ListNode{0, &ListNode{0, &ListNode{2, nil}}},
		},
	}
	for _, testCase := range testCases {
		if result := addTwoNumbers(testCase.leftList, testCase.righList); !reflect.DeepEqual(result, testCase.expectedList) {
			fmt.Print("Wanted: ")
			testCase.expectedList.PrintDownNodes()
			fmt.Print(" Returned: ")
			result.PrintDownNodes()
			t.Errorf("Expected result was incorrect.")
		}
	}
}
