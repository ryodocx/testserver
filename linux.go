//go:build linux

package main

import (
	"os"
	"syscall"
)

func init() {
	ignoreSignals = []os.Signal{syscall.SIGURG} // https://golang.hateblo.jp/entry/golang-signal-urgent-io-condition
}
