CURPATH=$(PWD)
TARGET_DIR=$(CURPATH)/bin

all build:
	CGO_ENABLED=0 GOOS=linux go build -trimpath -a -o $(TARGET_DIR)/cns-register cmd/register-cns-volume/main.go
.PHONY: all build

clean:
	rm -f $(TARGET_DIR)/playground
.PHONY: clean
