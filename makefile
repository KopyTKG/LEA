#Make file used for compilation and installation of LEA

.PHONY: all x86_64 arm64 clean

all: x86_64 arm64

x86_64:
	GOOS=linux GOARCH=amd64 go build -o build/lea.x64 main.go
	sha512sum build/lea.x64 >> build/lea.x64.sha512

arm64:
	GOOS=linux GOARCH=arm64 go build -o build/lea.arm64 main.go
	sha512sum build/lea.arm64 >> build/lea.arm64.sha512

clean:
	rm -f build/lea.x64 build/lea.arm64 build/lea.x64.sha512 build/lea.arm64.sha512
