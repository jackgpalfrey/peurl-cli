package utils

import (
	"fmt"
	"os"
)

func ExitWithError(message string, rawErr error) {
	fmt.Println(message)
	if os.Getenv("PEURL_DEBUG") == "enabled" {
		fmt.Println("Debug Error:", rawErr)
	}
	os.Exit(1)
}

func UsageError(usage string) {
	fmt.Println("Usage:", usage)
	os.Exit(2)
}
