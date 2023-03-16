build: *.go
	GOOS=windows GOARCH=amd64 go build -o windows/go-shell-api.exe -tags forceposix
	GOOS=darwin GOARCH=amd64 go build -o macos/go-shell-api -tags forceposix
	GOOS=linux GOARCH=amd64 go build -o linux-x64/go-shell-api -tags forceposix
	GOOS=linux GOARCH=arm64 go build -o linux-arm64/go-shell-api -tags forceposix
	GOOS=linux GOARCH=arm go build -o linux-arm/go-shell-api -tags forceposix
	GOOS=freebsd GOARCH=amd64 go build -o freebsd-x64/go-shell-api -tags forceposix
