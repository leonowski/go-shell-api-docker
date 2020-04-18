A utility to expose shell commands through a web API.  Copied from:
http://techfeast-hiranya.blogspot.com/2015/06/expose-any-shell-command-or-script-as.html

# Building locally:

Requirements:
golang

`go build`

# Dockerfile for automated builds for common platforms

After building, just copy the bins in /go-shell-api-bins.  Example:

`docker run -v "$(pwd)":/temp --rm -it noel-go-shell-test sh -c "cp -R /go-shell-api-bins/* /temp/"`

# Running

`go-shell-api -b some-binary-to-expose`

Examples:

`go-shell-api -b ls`

This will expose the command ls at http://ip-address:8080

An HTTP GET will give you the output of ls where the binary is running.  Use POST to add options to the command.  Example using curl:

`curl -X POST -d "-ls" some-ip-address:8080`
