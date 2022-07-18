package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// default config
var listenAddr = "127.0.0.1:8080"
var startupWait time.Duration = 0
var responseBody []byte = []byte("I'm a testserver")
var responseSleep time.Duration = 50 * time.Millisecond
var trapSignals []os.Signal = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
var gracePeriod time.Duration = 1 * time.Second

func init() {
	// override default config

	// LISTEN_ADDR
	if v := os.Getenv("LISTEN_ADDR"); v != "" {
		listenAddr = v
	}

	// STARTUP_WAIT
	if e := os.Getenv("STARTUP_WAIT"); e != "" {
		if v, err := time.ParseDuration(e); err != nil {
			log.Fatalln("invalid 'STARTUP_WAIT':", err)
		} else {
			startupWait = v
		}
	}

	// RESPONSE_BODY
	if v := os.Getenv("RESPONSE_BODY"); v != "" {
		responseBody = []byte(v)
	}

	// RESPONSE_SLEEP
	if e := os.Getenv("RESPONSE_SLEEP"); e != "" {
		if v, err := time.ParseDuration(e); err != nil {
			log.Fatalln("invalid 'RESPONSE_SLEEP':", err)
		} else {
			responseSleep = v
		}
	}

	// TRAP_SIGNALS
	if v := os.Getenv("TRAP_SIGNALS"); v != "" {
		trapSignals = []os.Signal{}
		if v == "0" {
			// disable graceful shutdown
		} else {
			// enable graceful shutdown
			for _, s := range strings.Split(v, ",") {
				i, err := strconv.Atoi(s)
				if err != nil {
					log.Fatalln("invalid 'TRAP_SIGNALS':", s)
				}
				trapSignals = append(trapSignals, syscall.Signal(i))
			}
		}
	}

	// GRACE_PERIOD
	if e := os.Getenv("GRACE_PERIOD"); e != "" {
		if v, err := time.ParseDuration(e); err != nil {
			log.Fatalln("invalid 'GRACE_PERIOD':", err)
		} else {
			gracePeriod = v
		}
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(responseSleep)
	w.Write(responseBody)
}

func main() {
	fmt.Println("pid =", os.Getpid())
	fmt.Println("############# Configuration #############")
	print := func(key string, val interface{}) { fmt.Printf("%-19s%v\n", key, val) }
	print("LISTEN_ADDR", listenAddr)
	print("STARTUP_WAIT", startupWait)
	print("RESPONSE_SLEEP", responseSleep)
	print("TRAP_SIGNALS", trapSignals)
	print("GRACE_PERIOD", gracePeriod)
	fmt.Println("#########################################")

	var srv http.Server
	srv.Addr = listenAddr
	http.HandleFunc("/", handler)

	idleConnsClosed := make(chan struct{})
	go func() {

		// signal monitoring
		for {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan)
			signal.Ignore(syscall.SIGURG) // https://golang.hateblo.jp/entry/golang-signal-urgent-io-condition
			recievedSignal := <-sigChan
			log.Println("signal recieved:", fmt.Sprintf("%d(%s)", recievedSignal, recievedSignal.String()))

			for _, s := range trapSignals {
				if recievedSignal == s {
					goto shutdown
				}
			}
		}

		// graceful shutdown
	shutdown:
		log.Println("waiting for shutdown:", gracePeriod)
		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()
		log.Println("shutting down...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// start
	if startupWait > 0 {
		log.Println("waiting for startup:", startupWait)
		time.Sleep(startupWait)
	}
	log.Println("start servering")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
