#Make file used for compilation and installation of LEA

install:
	go get -u
	go build
	sudo cp lea /usr/bin

build:
	go get -u
	go build

.PHONY: install build
