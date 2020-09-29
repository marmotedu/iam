# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for license
#
#
.PHONY: license.verify
license.verify: 
	@echo "===========> Verifying the boilerplate headers for all files"
	@$(GO) run $(ROOT_DIR)/tools/addlicense/addlicense.go --check -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,_output

.PHONY: license.add
license.add:
	@$(GO) run $(ROOT_DIR)/tools/addlicense/addlicense.go -v -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,_output
