# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for tools
#

DEP_TOOLS ?= swagger mockgen gotests gsemver golines go-junit-report git-chglog github-release coscmd go-mod-outdated golangci-lint protoc-gen-go cfssl addlicense goimports codegen
OTHER_TOOLS ?= depth go-callvis gothanks richgo rts

tools.install: $(addprefix tools.install., $(DEP_TOOLS), ${OTHER_TOOLS})
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.swagger
install.swagger:
	@$(GO) get -u github.com/go-swagger/go-swagger/cmd/swagger

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@golangci-lint completion bash > $(HOME)/.golangci-lint.bash
	@if ! grep -q .golangci-lint.bash $(HOME)/.bashrc; then echo "source \$$HOME/.golangci-lint.bash" >> $(HOME)/.bashrc; fi

.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) get -u github.com/jstemmer/go-junit-report

.PHONY: install.gsemver
install.gsemver:
	@$(GO) get -u github.com/arnaud-deprez/gsemver

.PHONY: install.git-chglog
install.git-chglog:
	@$(GO) get -u github.com/git-chglog/git-chglog/cmd/git-chglog

.PHONY: install.github-release
install.github-release:
	@$(GO) get -u github.com/github-release/github-release

.PHONY: install.coscmd
install.coscmd:
	@pip install coscmd

.PHONY: install.golines
install.golines:
	@$(GO) get -u github.com/segmentio/golines

.PHONY: install.go-mod-outdated
install.go-mod-outdated:
	@$(GO) get -u github.com/psampaz/go-mod-outdated

.PHONY: install.mockgen
install.mockgen:
	@$(GO) get -u github.com/golang/mock/mockgen

.PHONY: install.gotests
install.gotests:
	@$(GO) get -u github.com/cweill/gotests/...

.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	@$(GO) get -u github.com/golang/protobuf/protoc-gen-go

.PHONY: install.cfssl
install.cfssl:
	@$(ROOT_DIR)/scripts/install/install.sh iam::install::install_cfssl

.PHONY: install.addlicense
install.addlicense:
	@$(GO) get -u github.com/marmotedu/addlicense

.PHONY: install.goimports
install.goimports:
	@$(GO) get -u golang.org/x/tools/cmd/goimports

.PHONY: install.depth
install.depth:
	@$(GO) get -u github.com/KyleBanks/depth/cmd/depth

.PHONY: install.go-callvis
install.go-callvis:
	@$(GO) get -u github.com/ofabry/go-callvis

.PHONY: install.gothanks
install.gothanks:
	@$(GO) get -u github.com/psampaz/gothanks

.PHONY: install.richgo
install.richgo:
	@$(GO) get -u github.com/kyoh86/richgo

.PHONY: install.rts
install.rts:
	@$(GO) get -u github.com/galeone/rts/cmd/rts

.PHONY: install.codegen
install.codegen:
	@$(GO) install ${ROOT_DIR}/tools/codegen/codegen.go
