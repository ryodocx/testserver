//go:build linux

package main

import (
	"fmt"
	"os"
	"syscall"
)

func init() {
	ignoreSignals = []os.Signal{syscall.SIGURG}
	fmt.Println("hello")
}
