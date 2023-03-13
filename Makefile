.PHONY: build clean tool lint help
APP=glogin
DIR=phn_test-glogin-1
BINARY=${DIR}/${APP}
Version="v0.0.1"

ifeq ($(OS),Windows_NT)
 	PLATFORM="Windows"
else
 	ifeq ($(shell uname),Darwin)
  		PLATFORM="darwin"
 	else
  		PLATFORM="linux"
 	endif
endif


clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

all:
	@echo $(PLATFORM)

build:
	@go mod tidy
	@go build -ldflags "-X github.com/castbox/nirvana-kite/configs.version=${Version}" -o ${BINARY}

build_linux:
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build -o ./linux/${APP}

run:
	${BINARY}

help:
	@echo "本makefile 一共实现了以下几种命令模式"
	@echo "AllLibs \t\t默认的命令,含义为编译当前项目并输出可执行文件 $(target)"
	@echo "clean \t\t清理命令,删除由makefile带来的所有文件"
	@echo "clean_out \t\t清理中间文件命令,删除由编译带来的所有文件 *.o"