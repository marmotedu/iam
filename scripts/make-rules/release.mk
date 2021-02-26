# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for release
#
#

.PHONY: release.verify
release.verify: release.git-chglog.verify release.github-release.verify release.coscmd.verify

.PHONY: release.gsemver.verify                                      
release.gsemver.verify: 
ifeq (,$(shell which gsemver 2>/dev/null))    
	@$(MAKE) tools.install.gsemver
endif 

.PHONY: release.git-chglog.verify                                      
release.git-chglog.verify: 
ifeq (,$(shell which git-chglog 2>/dev/null))    
	@$(MAKE) tools.install.git-chglog
endif 

.PHONY: release.github-release.verify                                      
release.github-release.verify: 
ifeq (,$(shell which github-release 2>/dev/null))    
	@$(MAKE) tools.install.github-release
endif 

.PHONY: release.coscmd.verify                                      
release.coscmd.verify: 
ifeq (,$(shell which coscmd 2>/dev/null))    
	@$(MAKE) tools.install.coscmd
endif 

.PHONY: release.run
release.run: release.verify release.tag
	@scripts/release.sh

.PHONY: release.tag
release.tag: release.gsemver.verify
	@scripts/ensure_tag.sh
