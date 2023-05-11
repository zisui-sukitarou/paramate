BIN_NAME := paramate
BIN_DIR := ./bin
X_BIN_DIR := $(BIN_DIR)/goxz
VERSION := $$(make -s app-version)

GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: all
all: build

# 実行環境に適したバイナリを生成する
.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) main.go

# 複数のプラットフォームを対象にそれぞれの環境に適したバイナリを生成しzip化する
.PHONY: x-build
x-build: $(GOBIN)/goxz
	goxz -d $(X_BIN_DIR) -n $(BIN_NAME) .

# 生成したバイナリをGitHubのReleaseにアップロードする (x-buildで生成したものを対象にしています)
.PHONY: upload-binary
upload-binary: $(GOBIN)/ghr
	ghr "v$(VERSION)" $(X_BIN_DIR)

# アプリのバージョンを出力する
.PHONY: app-version 
app-version: $(GOBIN)/gobump
	@gobump show -r .

# 以下、上記のターゲットにて使用するツールが導入されていなかった場合に`go install`で導入を行う
$(GOBIN)/goxz:
	@go install github.com/Songmu/goxz/cmd/goxz@latest

$(GOBIN)/ghr:
	@go install github.com/tcnksm/ghr@latest

$(GOBIN)/gobump:
	@go install github.com/x-motemen/gobump/cmd/gobump@master