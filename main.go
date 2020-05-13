package main

import (
	"fmt"
	"regexp"
)

func main() {

	fmt.Println("Phone Number Normalizer")

}

func normalizer(phone string) string {

	re := regexp.MustCompile("[^0-9]")

	return re.ReplaceAllString(phone, "")
}

// func normalizer(phone string) string {

// 	var buf bytes.Buffer

// 	for _, ch := range phone {

// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}

// 	}
// 	return buf.String()
// }
