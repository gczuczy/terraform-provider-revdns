NAME=revdns
HOSTNAME=github.com
NAMESPACE=gczuczy
VERSION=0.1.0
BINARY=terraform-provider-${NAME}_v${VERSION}

uname_s := $(shell uname -s)
uname_m := $(shell uname -m)
OS_ARCH.FreeBSD.amd64 := freebsd_amd64
OS_ARCH.Linux.x86_64 := linux_amd64
OS_ARCH.Darwin.x86_64 := linux_amd64
OS_ARCH = $(OS_ARCH.$(uname_s).$(uname_m))
$(info OS_ARCH=$(OS_ARCH))
INSTALLDIR=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

default: build-dev

build:
	go build -o ${BINARY}

build-dev:
	mkdir -p ~/.terraform.d/plugins/
	go build -o ~/.terraform.d/plugins/${BINARY}

install: build
	mkdir -p ${INSTALLDIR}
	mv ${BINARY} ${INSTALLDIR}/
