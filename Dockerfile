FROM golang:1.20 as builder
ADD main.go Makefile /go-shell-api-docker/
WORKDIR /go-shell-api-docker/
RUN go env -w GO111MODULE=off && make build && rm main.go Makefile
FROM alpine
COPY --from=builder /go-shell-api-docker /go-shell-api-bins/
ENTRYPOINT ["echo","This container does nothing.  Copy the files you want from /go-shell-api-bins"]
