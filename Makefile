.PHONY: all

# Build native go app
dm2: *.go
	mkdir -p dist
	go build -o dist/$@ .

# Build for arm64 Mac
dm2_mac_arm64: *.go
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build -o dist/$@ .

# Build for x64 linux
dm2_linux_x64: *.go
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/$@ .

all: dm2 dm2_linux_x64 dm2_mac_arm64