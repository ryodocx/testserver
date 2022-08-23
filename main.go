package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// build info
var (
	version string
)

// default config
var listenAddr = "127.0.0.1:8080"
var startupWait time.Duration = 0
var responseBody []byte = []byte("I'm a testserver")
var responseSleep time.Duration = 50 * time.Millisecond
var trapSignals []os.Signal = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
var gracePeriodBeforeShutdown time.Duration = 1 * time.Second
var gracePeriodDuringShutdown time.Duration = 0
var accessLog bool = false

var ignoreSignals []os.Signal

func init() {
	// override default config

	getEnvDuration := func(envName string, defaultValue time.Duration) time.Duration {
		e := os.Getenv(envName)
		if e == "" {
			return defaultValue
		}
		v, err := time.ParseDuration(e)
		if err != nil {
			log.Fatalln("invalid ", envName, ":", err)
		}
		return v
	}

	// LISTEN_ADDR
	if v := os.Getenv("LISTEN_ADDR"); v != "" {
		listenAddr = v
	}

	// STARTUP_WAIT
	startupWait = getEnvDuration("STARTUP_WAIT", startupWait)

	// RESPONSE_BODY
	if v := os.Getenv("RESPONSE_BODY"); v != "" {
		responseBody = []byte(v)
	}

	// RESPONSE_SLEEP
	responseSleep = getEnvDuration("RESPONSE_SLEEP", responseSleep)

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

	// GRACE_PERIOD_BEFORE_SHUTDOWN
	gracePeriodBeforeShutdown = getEnvDuration("GRACE_PERIOD_BEFORE_SHUTDOWN", gracePeriodBeforeShutdown)

	// GRACE_PERIOD_DURING_SHUTDOWN
	gracePeriodDuringShutdown = getEnvDuration("GRACE_PERIOD_DURING_SHUTDOWN", gracePeriodDuringShutdown)

	// ACCESS_LOG
	if e := os.Getenv("ACCESS_LOG"); e == "true" {
		accessLog = true
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	if accessLog {
		log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.RequestURI)
	}
	time.Sleep(responseSleep)
	_, _ = w.Write(responseBody)
}

func main() {
	print := func(key string, val interface{}) { fmt.Printf("%-29s%v\n", key, val) }
	i, _ := debug.ReadBuildInfo()
	if version == "" {
		if i.Main.Version != "(devel)" {
			version = i.Main.Version
		} else {
			version = "unknown"
		}
	}
	fmt.Println("###################### Info #######################")
	fmt.Println(i.Path)
	print("Version", version)
	print("GoVersion", i.GoVersion)
	print("PID", os.Getpid())
	fmt.Println("################## Configuration ##################")
	print("LISTEN_ADDR", listenAddr)
	print("STARTUP_WAIT", startupWait)
	print("RESPONSE_SLEEP", responseSleep)
	print("TRAP_SIGNALS", trapSignals)
	print("GRACE_PERIOD_BEFORE_SHUTDOWN", gracePeriodBeforeShutdown)
	print("GRACE_PERIOD_DURING_SHUTDOWN", gracePeriodDuringShutdown)
	fmt.Println("###################################################")

	var srv http.Server
	srv.Addr = listenAddr
	http.HandleFunc("/", handler)

	idleConnsClosed := make(chan struct{})
	go func() {

		// signal monitoring
		for {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan)
			if len(ignoreSignals) > 0 {
				signal.Ignore(ignoreSignals...)
			}
			receivedSignal := <-sigChan
			log.Println("signal received:", fmt.Sprintf("%d(%s)", receivedSignal, receivedSignal.String()))

			for _, s := range trapSignals {
				if receivedSignal == s {
					goto shutdown
				}
			}
		}

		// graceful shutdown
	shutdown:
		log.Println("waiting for shutdown:", gracePeriodBeforeShutdown)
		time.Sleep(gracePeriodBeforeShutdown)

		ctx := context.Background()
		if gracePeriodDuringShutdown == 0 {
			log.Println("shutting down... (grace period = unlimited)")
		} else {
			log.Printf("shutting down... (grace period = %v)\n", gracePeriodDuringShutdown)
			c, cancel := context.WithTimeout(context.Background(), gracePeriodDuringShutdown)
			defer cancel()
			ctx = c
		}
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
	log.Println("server start")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
