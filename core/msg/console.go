package msg

import (
    "fmt"
    "os"
)

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
    return &ConsoleLogger{}
}

func (ConsoleLogger) Info(msg string) {
    fmt.Println("\033[34m[i]\033[0m", msg)
}

func (ConsoleLogger) Success(msg string) {
    fmt.Println("\033[32m[✓]\033[0m", msg)
}

func (ConsoleLogger) Warn(msg string) {
    fmt.Println("\033[33m[!]\033[0m", msg)
}

func (ConsoleLogger) Error(msg string) {
    fmt.Println("\033[31m[x]\033[0m", msg)
}

func (ConsoleLogger) Fail(err error) {
    if err != nil {
        fmt.Println("\033[31m[x]\033[0m", err.Error())
        os.Exit(1)
    }
}
