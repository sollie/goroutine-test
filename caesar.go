package main

import (
	"unicode"
)

func init() {
	registerFunction("caesar", caesar)
}

func caesar(input string, shift int) string {
	runes := []rune(input)
	shifted := make([]rune, len(runes))

	for i, char := range runes {
		if unicode.IsLetter(char) {
			var base rune
			if unicode.IsUpper(char) {
				base = 'A'
			} else {
				base = 'a'
			}
			shifted[i] = (char-base+rune(shift))%26 + base
		} else {
			shifted[i] = char
		}
	}

	return string(shifted)
}
