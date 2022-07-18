FROM golang:1.18.4-alpine
ENV CGO_ENABLED=0
WORKDIR /
COPY . .
RUN go build -o /testserver *.go

FROM alpine:3.16.0
ENV LISTEN_ADDR=0.0.0.0:8080
COPY --from=0 /testserver .
ENTRYPOINT [ "/testserver" ]
