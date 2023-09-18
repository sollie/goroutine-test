package main

import (
	"strings"
)

func init() {
	registerFunction("uppercase", uppercase)
}

func uppercase(s string) string {
	return strings.ToUpper(s)
}
