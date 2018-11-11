// Objective: Find the length of the longest substring without
// repeating characters in a given string.

package main

import (
	"fmt"

	"github.com/pkg/profile"
)

func lengthOfLongestSubstring(s string) int {

	if s == "" {
		return 0
	}

	if len(s) == 1 {
		return 1
	}

	maxStrLen := 1
	currentStrLen := 1
	seen := make(map[byte]bool)
	hit := false
	for l, r := 0, 0; r < len(s); {
		if _, hit = seen[s[r]]; !hit {
			seen[s[r]] = true
			r++
			if currentStrLen = len(s[l:r]); currentStrLen > maxStrLen {
				maxStrLen = currentStrLen
			}
		} else {
			delete(seen, s[l])
			l++
		}

	}

	return maxStrLen
}

func main() {
	defer profile.Start().Stop()
	myString := "applebees" // 4 ('pleb') should be the answer
	fmt.Println("Length of maximum substirng without characters repeating is:", lengthOfLongestSubstring(myString))

}
