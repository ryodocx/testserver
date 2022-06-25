# testserver
HTTP Server for graceful-shutdown testing

## Usage

```sh
# go install
go install github.com/ryodocx/testserver@latest

# docker
docker run --rm -it -p 8080:8080 ghcr.io/ryodocx/testserver
```

## Environment variables

| env            | default                  | example                                                                                                     | description                                                           |
|----------------|--------------------------|-------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------|
| LISTEN_ADDR    | `0.0.0.0:8080`           | `127.0.0.1:8080`                                                                                            | Listen address                                                        |
| RESPONSE_BODY  | `I'm a testserver`       | `hello world`                                                                                               | HTTP response body                                                    |
| RESPONSE_SLEEP | `50ms`                   | `0` (without sleep) <br> `200ms` `5s` `0.01h`                                                               | Sleep time during HTTP response                                       |
| TRAP_SIGNALS   | `[interrupt terminated]` | `0` (disable graceful shutdown) <br> `1,2,15` (enable graceful shutdown for SIGHUP/SIGINT/SIGTERM at Linux) | Trapped Signals for graceful shutdown                                 |
| GRACE_PERIOD   | `1s`                     | `0` (no wait) <br> `5s` `1m`                                                                                | Grace period before starting shutdown (ignored when `TRAP_SIGNALS=0`) |
