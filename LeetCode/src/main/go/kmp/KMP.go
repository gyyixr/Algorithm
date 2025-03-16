package main

import "fmt"

func kmpSearch(text, pattern string) []int {
	m := len(pattern)
	n := len(text)
	lps := computeLPSArray(pattern)

	var matches []int
	i, j := 0, 0
	for i < n {
		if pattern[j] == text[i] {
			i++
			j++
		}
		if j == m {
			matches = append(matches, i-j)
			j = lps[j-1]
		} else if i < n && pattern[j] != text[i] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}
	return matches
}

func computeLPSArray(pattern string) []int {
	m := len(pattern)
	lps := make([]int, m)
	len := 0
	i := 1
	for i < m {
		if pattern[i] == pattern[len] {
			len++
			lps[i] = len
			i++
		} else {
			if len != 0 {
				len = lps[len-1]
			} else {
				lps[i] = 0
				i++
			}
		}
	}
	return lps
}

func main() {
	text := "ABABDABACDABABCABAB"
	pattern := "ABABCABAB"
	matches := kmpSearch(text, pattern)
	if len(matches) > 0 {
		fmt.Printf("Pattern found at index/indices: %v\n", matches)
	} else {
		fmt.Println("Pattern not found.")
	}
}
