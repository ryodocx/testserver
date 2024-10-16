FROM golang:1.22.3-alpine
RUN apk add git
ENV CGO_ENABLED=0
WORKDIR /
COPY . .
RUN go install -ldflags "-X main.version=$(git describe --tags)"

FROM alpine:3.19.1
ENV LISTEN_ADDR=0.0.0.0:8080
COPY --from=0 /go/bin/testserver /usr/local/bin/
CMD [ "testserver" ]
