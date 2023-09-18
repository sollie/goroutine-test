package main

func init() {
	registerFunction("reverse", reverse)
}

func reverse(s string) string {
	data := []rune(s)
	result := []rune{}
	for i := len(data) - 1; i >= 0; i-- {
		result = append(result, data[i])
	}
	return string(result)
}
