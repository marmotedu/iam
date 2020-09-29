# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for release
#
#

release.verify: release.git-chglog.verify release.github-release.verify release.coscmd.verify

.PHONY: release.gsemver.verify                                      
release.gsemver.verify: 
ifeq (,$(shell which gsemver 2>/dev/null))    
	@echo "===========> Installing gsemver"
	@$(GO) get -u github.com/arnaud-deprez/gsemver
endif 

.PHONY: release.git-chglog.verify                                      
release.git-chglog.verify: 
ifeq (,$(shell which git-chglog 2>/dev/null))    
	@echo "===========> Installing git-chglog"
	@$(GO) get -u github.com/git-chglog/git-chglog/cmd/git-chglog
endif 

.PHONY: release.github-release.verify                                      
release.github-release.verify: 
ifeq (,$(shell which github-release 2>/dev/null))    
	@echo "===========> Installing github-release"
	@$(GO) get -u github.com/github-release/github-release
endif 

.PHONY: release.coscmd.verify                                      
release.coscmd.verify: 
ifeq (,$(shell which coscmd 2>/dev/null))    
	@echo "===========> Installing coscmd"
	@pip install coscmd
endif 

.PHONY: release.run
release.run: release.verify release.tag
	@scripts/release.sh

.PHONY: release.tag
release.tag: release.gsemver.verify
	@scripts/ensure_tag.sh
