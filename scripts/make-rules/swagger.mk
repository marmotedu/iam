# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.    
# Use of this source code is governed by a MIT style    
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for swagger
#

.PHONY: swagger.verify
swagger.verify: 
ifeq (,$(shell which swagger 2>/dev/null))
	@echo "===========> Installing go-swagger"
	@$(GO) get -u github.com/go-swagger/go-swagger/cmd/swagger
endif

.PHONY: swagger.run
swagger.run: swagger.verify
	@echo "===========> Generating swagger API docs"
	@swagger generate spec --scan-models -w $(ROOT_DIR)/cmd/genswaggertypedocs -o $(ROOT_DIR)/api/swagger/swagger.yaml

.PHONY: swagger.serve
swagger.serve: swagger.verify
	@swagger serve -F=redoc --no-open --port 36666 $(ROOT_DIR)/api/swagger/swagger.yaml
