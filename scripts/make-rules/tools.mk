# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for tools 
#

DEP_TOOLS ?= swagger mockgen gotests gsemver golines go-junit-report git-chglog github-release coscmd go-mod-outdated golangci-lint protoc-gen-go

tools.install: $(addprefix tools.install., $(DEP_TOOLS))

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
