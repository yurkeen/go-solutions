package main

import "fmt"

// Objective: Given two non-empty linked lists representing two non-negative integers,
// add the two numbers and return it as a linked list.
// The digits are stored in reverse order and each of their nodes contain a single digit.
//
// Assumptions: the two numbers do not contain any leading zero, except the number 0 itself.
// Example:

// Input: (4 -> 5 -> 6) + (1 -> 2 -> 3)
// Output: 9 -> 7 -> 5
// Explanation: 654 + 321 = 975.

// ListNode defines a singly-linked list
type ListNode struct {
	Val  int
	Next *ListNode
}

// PrintNodes traverses the list down printing all values
func (l *ListNode) PrintNodes() {
	for node := l; node != nil; node = node.Next {
		fmt.Printf("[%d]-> ", node.Val)
	}
}

func (l *ListNode) listLen() int {
	result := 0
	for node := l; node != nil; node = node.Next {
		result++
	}
	return result
}

// findLast returns pointer to the last node in the linked list
func (l *ListNode) findLast() *ListNode {
	for node := l; ; {
		if node.Next == nil {
			return node
		}
		node = node.Next
	}
}

// appendNode appends a new node with given value to last position in the linked list and returns back pointer to a new node
func (l *ListNode) appendNode(v int) *ListNode {
	node := l.findLast()
	newNode := &ListNode{Val: v, Next: nil}
	node.Next = newNode
	return newNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	// We want to add to the largest (longer) value
	longer, shorter := l1, l2
	if l1.listLen() < l2.listLen() {
		longer, shorter = l2, l1
	}
	carry := 0 // Handles digit overflow
	for cursor1, cursor2 := longer, shorter; cursor1 != nil; cursor1 = cursor1.Next {
		// Cursors point to node in respective lists.
		// We iterate over the longer list.
		cursor1.Val += carry

		valSum := cursor1.Val
		if cursor2 != nil {
			valSum += cursor2.Val
			cursor2 = cursor2.Next
		}
		carry = valSum / 10
		cursor1.Val = valSum % 10

		if carry != 0 && cursor1.Next == nil {
			cursor1.appendNode(carry)
			break
		}
	}
	return longer
}

func main() {
	firstNumber := &ListNode{Val: 4, Next: nil}
	firstNumber.appendNode(5)
	firstNumber.appendNode(6)

	secondNumber := &ListNode{Val: 1, Next: nil}
	secondNumber.appendNode(2)
	secondNumber.appendNode(3)

	result := addTwoNumbers(firstNumber, secondNumber)
	result.PrintNodes()
}
