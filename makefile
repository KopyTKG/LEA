#Make file used for compilation and installation of LEA

install:
	go get

build:
	go build

sysinstall:
	sudo cp lea /usr/bin/

all: install build sysinstall

.PHONY: install build sysinstall all 
