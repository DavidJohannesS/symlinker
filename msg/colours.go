package msg

import (
	"fmt"
)
func Info(msg string) {
    fmt.Println("\033[34m[i]\033[0m", msg)
}

func Success(msg string) {
    fmt.Println("\033[32m[✓]\033[0m", msg)
}

func Warn(msg string) {
    fmt.Println("\033[33m[!]\033[0m", msg)
}

func Error(msg string) {
    fmt.Println("\033[31m[x]\033[0m", msg)
}
