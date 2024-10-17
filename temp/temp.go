package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "john_tah"
	idx := strings.LastIndex(str, "_")
	if idx != -1 {
		substr := str[idx+1:]
		fmt.Println(substr) // Output: tah
	}
}
