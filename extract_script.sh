#!/bin/bash
docker run -v "$(pwd)":/temp --rm -it noel-go-shell-test sh -c "cp -R /go-shell-api-bins/* /temp/"
