// Objective: Given an array of integers, return indices of the two numbers such that they add up to a specific target.
// Assumptions: each input would have exactly one solution, same element cannot be used twice.
//
// Example:
// Given nums = [2, 7, 11, 15], target = 9,
// Because nums[0] + nums[1] = 2 + 7 = 9,
// return [0, 1].
package main

import (
	"fmt"
)

func findAddend(numbers []int, sum int) []int {
	// map allows us for fast lookups of exiting
	// elements by value
	lookup := make(map[int]int)
	for idx1, v := range numbers {
		if idx2, ok := lookup[sum-numbers[idx1]]; ok {
			return []int{idx2, idx1}
		}
		lookup[v] = idx1
	}
	return nil
}

func main() {
	nums := []int{10, 4, 1, -20, 0, 40, -2, 7, 11, 15}
	target := 9
	result := findAddend(nums, target)
	fmt.Println("Result is:", result)
}
