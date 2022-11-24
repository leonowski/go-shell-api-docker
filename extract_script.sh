#!/bin/bash
docker run -v "$(pwd)":/temp --entrypoint=/bin/cp --rm -it leonowski/go-shell-api-docker -R /go-shell-api-bins /temp/
