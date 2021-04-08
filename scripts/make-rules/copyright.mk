# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for copyright
#
#
.PHONY: copyright.verify
copyright.verify: tools.verify.addlicense
	@echo "===========> Verifying the boilerplate headers for all files"
	@addlicense --check -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,vendor,_output

.PHONY: copyright.add
copyright.add: tools.verify.addlicense
	@addlicense -v -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,vendor,_output
