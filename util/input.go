package util

import "fmt"

func Read(message string) string {
	var input string

	fmt.Print(message)
	_, err := fmt.Scanln(&input)

	if err != nil {
		fmt.Println("Error reading input:", err)
		Read(message)
	}

	return input
}
