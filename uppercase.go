package main

import (
	"strings"
)

func init() {
	registerFunction("uppercase", uppercase)

	workers = append(workers, Worker{
		Function: "uppercase",
		Args: []interface{}{
			"Me gustan los tacos",
		},
		Sleep: 8,
	},
	)
}

func uppercase(s string) string {
	return strings.ToUpper(s)
}
