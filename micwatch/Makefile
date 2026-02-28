APP_NAME := micwatch
LABEL := com.andrey.micwatch

BIN_DIR := $(HOME)/bin
LAUNCH_AGENTS_DIR := $(HOME)/Library/LaunchAgents

BIN_PATH := $(BIN_DIR)/$(APP_NAME)
PLIST_PATH := $(LAUNCH_AGENTS_DIR)/$(LABEL).plist

PLIST_TEMPLATE := micwatch.plist.in

OUT_LOG := /tmp/$(APP_NAME).out.log
ERR_LOG := /tmp/$(APP_NAME).err.log

BUILD_DIR := build
ARM_BIN := $(BUILD_DIR)/$(APP_NAME)-arm64
X64_BIN := $(BUILD_DIR)/$(APP_NAME)-x86_64
UNIV_BIN := $(BUILD_DIR)/$(APP_NAME)

MACOSX_DEPLOYMENT_TARGET ?= 12.0

HOST ?= 127.0.0.1

UID := $(shell id -u)
GUI_DOMAIN := gui/$(UID)

.PHONY: build clean install start stop restart status uninstall plist

build:
	mkdir -p $(BUILD_DIR)

	MACOSX_DEPLOYMENT_TARGET=$(MACOSX_DEPLOYMENT_TARGET) \
	swiftc -O -target arm64-apple-macos$(MACOSX_DEPLOYMENT_TARGET) \
	-framework CoreAudio -framework AudioToolbox \
	main.swift -o $(ARM_BIN)

	MACOSX_DEPLOYMENT_TARGET=$(MACOSX_DEPLOYMENT_TARGET) \
	swiftc -O -target x86_64-apple-macos$(MACOSX_DEPLOYMENT_TARGET) \
	-framework CoreAudio -framework AudioToolbox \
	main.swift -o $(X64_BIN)

	lipo -create -output $(UNIV_BIN) $(ARM_BIN) $(X64_BIN)
	file $(UNIV_BIN)

install: build
	mkdir -p $(BIN_DIR)
	cp -f $(UNIV_BIN) $(BIN_PATH)
	chmod +x $(BIN_PATH)
	@echo "Binary installed to $(BIN_PATH)"
	@echo "Now run: make start HOST=192.168.1.X"

plist:
	mkdir -p $(LAUNCH_AGENTS_DIR)
	@test -f $(PLIST_TEMPLATE) || (echo "Missing template: $(PLIST_TEMPLATE)"; exit 1)
	@echo "Writing $(PLIST_PATH) from template (HOST=$(HOST))"

	sed \
	  -e 's|__LABEL__|$(LABEL)|g' \
	  -e 's|__BIN_PATH__|$(BIN_PATH)|g' \
	  -e 's|__HOST__|$(HOST)|g' \
	  -e 's|__OUT_LOG__|$(OUT_LOG)|g' \
	  -e 's|__ERR_LOG__|$(ERR_LOG)|g' \
	  $(PLIST_TEMPLATE) > $(PLIST_PATH)

start: plist
	@echo "Starting LaunchAgent with HOST=$(HOST)"
	launchctl bootstrap $(GUI_DOMAIN) $(PLIST_PATH) 2>/dev/null || true
	launchctl enable $(GUI_DOMAIN)/$(LABEL) 2>/dev/null || true
	@echo "Logs:"
	@echo "  tail -f $(OUT_LOG)"
	@echo "  tail -f $(ERR_LOG)"

stop:
	launchctl bootout $(GUI_DOMAIN) $(PLIST_PATH) 2>/dev/null || true

restart: stop start

status:
	launchctl print $(GUI_DOMAIN)/$(LABEL) 2>/dev/null | head -n 60 || echo "Not loaded"

uninstall: stop
	rm -f $(PLIST_PATH)
	rm -f $(BIN_PATH)
	@echo "Uninstalled"

clean:
	rm -rf $(BUILD_DIR)