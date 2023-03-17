FROM golang:1.20 as builder
# Default allows all binaries/commands to run!!  Use a comma separated list of binaries at docker build time using ARG.  Example:  --build-arg "DOCKER_ALLOWED_COMMANDS=ls,ps,df"
ARG DOCKER_ALLOWED_COMMANDS=
ADD main.go Makefile /go-shell-api-docker/
WORKDIR /go-shell-api-docker/
RUN go env -w GO111MODULE=off && make build ALLOWED_COMMANDS=${DOCKER_ALLOWED_COMMANDS} && rm main.go Makefile
FROM alpine
COPY --from=builder /go-shell-api-docker /go-shell-api-bins/
ENTRYPOINT ["echo","This container does nothing.  Copy the files you want from /go-shell-api-bins"]
