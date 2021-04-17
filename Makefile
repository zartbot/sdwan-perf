GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
LDFLAGS = -s -w
BINARY_NAME=sdwan-perf
BINARY_PATH=./build/

all: build
build:  linux mac win linuxarm rpi
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_PATH)
	mkdir -p $(BINARY_PATH)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

linux:  $(info >> Starting build linux x86 based package...)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BINARY_PATH)/$(BINARY_NAME)_linux -v

linuxarm: $(info >> Starting build linux arm based package...)
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BINARY_PATH)/$(BINARY_NAME)_linux_arm -v

rpi: $(info >> Starting build raspberry pi based package...)
	GOOS=linux GOARCH=arm $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BINARY_PATH)/$(BINARY_NAME)_rpi -v

mac:  $(info >> Starting build mac package...)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BINARY_PATH)/$(BINARY_NAME)_mac -v

win:  $(info >> Starting build windows package...)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BINARY_PATH)/$(BINARY_NAME)_win.exe -v


