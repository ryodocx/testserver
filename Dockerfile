FROM golang:1.18.3-alpine
COPY *.go .
ENV CGO_ENABLED=0
RUN go build -o server *.go

FROM scratch
ENV LISTEN_ADDR=0.0.0.0:8080
COPY --from=0 /go/server .
CMD [ "/server" ]
