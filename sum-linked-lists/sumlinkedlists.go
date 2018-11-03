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

// PrintDownNodes traverses the list down printing all values
func (l *ListNode) PrintDownNodes() {
	for node := l; node != nil; node = node.Next {
		fmt.Printf("[%d]", node.Val)
		if node.Next != nil {
			fmt.Printf("->")
		}
	}
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	result := &ListNode{}

	cursor1 := l1
	cursor2 := l2

	// Iterate until both cursors are nill and carry is empty
	for sum, carry, node := 0, 0, result; ; {
		sum = carry
		if cursor1 != nil {
			sum += cursor1.Val
			cursor1 = cursor1.Next
		}

		if cursor2 != nil {
			sum += cursor2.Val
			cursor2 = cursor2.Next
		}

		carry = sum / 10
		node.Val = sum % 10

		if cursor1 == nil && cursor2 == nil && carry == 0 {
			return result
		}

		node.Next = &ListNode{carry, nil}
		node = node.Next
	}
}

func main() {
	firstNumber := &ListNode{Val: 4, Next: nil}
	firstNumber.appendNode(5)
	firstNumber.appendNode(6)

	secondNumber := &ListNode{Val: 1, Next: nil}
	secondNumber.appendNode(2)
	secondNumber.appendNode(3)

	result := addTwoNumbers(firstNumber, secondNumber)
	result.PrintDownNodes()
}
