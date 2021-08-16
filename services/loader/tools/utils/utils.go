package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals
	fmt.Println("Got signal for exiting", sig)
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func findIndex(s []string, item string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == item {
			return i
		}
	}
	return -1
}

func RemoveByItem(s []string, item string) []string {
	i := findIndex(s, item)

	if i == -1 {
		return s
	}

	return append(s[:i], s[i+1:]...)
}
