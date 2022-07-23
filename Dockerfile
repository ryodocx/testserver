FROM golang:1.18.4-alpine
RUN apk add git
ENV CGO_ENABLED=0
WORKDIR /
COPY . .
RUN go build \
	-ldflags \
	" \
        -X main.version=$(git describe --tag --abbrev=0) \
        -X main.revision=$(git rev-parse HEAD) \
        -X main.changed=$(git status -s | wc -l) \
    " \
	.

FROM alpine:3.16.0
ENV LISTEN_ADDR=0.0.0.0:8080
COPY --from=0 /testserver .
ENTRYPOINT [ "/testserver" ]
