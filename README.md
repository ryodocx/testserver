# testserver

HTTP Server for graceful-shutdown testing

## Usage

```sh
# go install
go install github.com/ryodocx/testserver@latest

# docker
docker run --rm -it -p 8080:8080 ghcr.io/ryodocx/testserver:v1
```

## Environment variables

| env                          | default                  | example                                                                                                     | description                                                           |
|------------------------------|--------------------------|-------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------|
| LISTEN_ADDR                  | `0.0.0.0:8080`           | `127.0.0.1:8080`                                                                                            | Listen address                                                        |
| STARTUP_WAIT                 | `0s`                     | `3s`                                                                                                        | Waiting time before start serving                                     |
| RESPONSE_BODY                | `I'm a testserver`       | `hello world`                                                                                               | HTTP response body                                                    |
| RESPONSE_SLEEP               | `50ms`                   | `0` (without sleep) <br> `200ms` `5s` `0.01h`                                                               | Wait time at HTTP response                                            |
| TRAP_SIGNALS                 | `[interrupt terminated]` | `0` (disable graceful shutdown) <br> `1,2,15` (enable graceful shutdown for SIGHUP/SIGINT/SIGTERM at Linux) | Trapped Signals for graceful shutdown                                 |
| GRACE_PERIOD_BEFORE_SHUTDOWN | `1s`                     | `0` (no wait) <br> `5s` `1m`                                                                                | Grace period before starting shutdown (ignored when `TRAP_SIGNALS=0`) |
| GRACE_PERIOD_DURING_SHUTDOWN | `0` (unlimited)          | `0` (unlimited) <br> `5s` `1m`                                                                              | Grace period during shutdown          (ignored when `TRAP_SIGNALS=0`) |
| ACCESS_LOG                   | `false`                  | `true`                                                                                                      | If true, enable access logging                                        |
