A utility to expose shell commands through a web API.  Copied from:
http://techfeast-hiranya.blogspot.com/2015/06/expose-any-shell-command-or-script-as.html

# Building locally

Requirements:
golang

`go build`

The Makefile builds for multiple platforms.

`go env -w GO111MODULE=off && make build`

You can now restrict the builds to only allow certain shell commands:

`go env -w GO111MODULE=off && make build make build ALLOWED_COMMANDS=ls,df,ps,free`


# Using Docker to build (recommended as you don't need GO installed locally.  Just need Docker)

Building in Docker
`docker build whatever-you-want-to-call-this-thing .`

It is now possible to restrict the shell commands allowed to be run through this tool.  You can use a Docker build arg to achieve this:

`docker build -t latest --build-arg "DOCKER_ALLOWED_COMMANDS=ls,ps" .`

This will output binaries in the container that only allow the commands "ls" and ps".  So, you just need a commma separated list of binaries.

After building, just copy the bins in /go-shell-api-bins.  Example:

`docker run -v "$(pwd)":/temp --rm -it whatever-you-want-call-this-thing sh -c "cp -R /go-shell-api-bins/* /temp/"`

Pre-made container available on Dockerhub (this container allows all binaries!):

leonowski/go-shell-api-docker

# Running

`go-shell-api -b some-binary-to-expose` (some-binary-to-expose is any shell command)

Examples:

`go-shell-api -b ls`

This will expose the command ls at http://ip-address:8080

An HTTP GET will give you the output of ls where the go-shell-api binary is running.  Use POST to add options to the command.  Example using curl:

`curl -X POST -d "-ls" some-ip-address:8080`

This will result in an output of the command `ls -ls` in the directory where go-shell-api is running.

You can specify HTTPS and basic auth with these options:

`./go-shell-api -b ls -https -u username -pw somepassword`

If HTTPS is used, a self-signed cert is used.  However, a cert and key can also be specified if desired.  See full options below.

Basic auth can be required as well using the `-u` and `-pw` options

Full options available with -help:

```
  -b string
    	Path to the executable binary
  -cert string
    	Path to the server certificate file (requires -https flag)
  -https
    	Serve via HTTPS using a self-signed certificate or an optional custom certificate
  -key string
    	Path to the server key file (requires -https flag)
  -l string
    	Address to listen on (default "0.0.0.0")
  -p int
    	HTTP port to listen on (default 8080)
  -pw string
    	Basic authentication password
  -u string
    	Basic authentication username
```
