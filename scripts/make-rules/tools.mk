# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for tools 
#

DEP_TOOLS ?= swagger golangci-lint go-junit-report gsemver git-chglog github-release coscmd golines

tools.install: $(addprefix tools.install., $(DEP_TOOLS))

tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) --no-print-directory install.$*

tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) --no-print-directory tools.install.$*; fi

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
