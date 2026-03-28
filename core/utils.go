package core

import (
	"fmt"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func PrintSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

func PrintError(message string) {
	fmt.Printf("❌ %s\n", message)
}

func PrintInfo(message string) {
	fmt.Printf("📡 %s\n", message)
}
