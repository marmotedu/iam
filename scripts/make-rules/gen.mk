# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for generate necessary files
#

PROTOC_INC_PATH=$(dir $(shell which protoc 2>/dev/null))/../include    
API_DEPS=pkg/model/apiserver/v1/cache.proto             
API_DEPSRCS=$(API_DEPS:.proto=.pb.go) 

.PHONY: gen.run
#gen.run: gen.errcode gen.docgo
gen.run: gen.clean gen.proto gen.errcode gen.docgo.doc

.PHONY: gen.errcode
gen.errcode: gen.errcode.code gen.errcode.doc

.PHONY: gen.errcode.code
gen.errcode.code: tools.verify.codegen
	@echo "===========> Generating iam error code go source files"
	@codegen -type=int ${ROOT_DIR}/internal/pkg/code

.PHONY: gen.errcode.doc
gen.errcode.doc: tools.verify.codegen
	@echo "===========> Generating error code markdown documentation"
	@codegen -type=int -doc \
		-output ${ROOT_DIR}/docs/guide/zh-CN/api/error_code_generated.md ${ROOT_DIR}/internal/pkg/code

.PHONY: gen.ca.%
gen.ca.%:
	$(eval CA := $(word 1,$(subst ., ,$*)))
	@echo "===========> Generating CA files for $(CA)"
	@${ROOT_DIR}/scripts/gencerts.sh generate-iam-cert $(OUTPUT_DIR)/cert $(CA)

.PHONY: gen.ca
gen.ca: $(addprefix gen.ca., $(CERTIFICATES))

.PHONY: gen.docgo.doc
gen.docgo.doc:
	@echo "===========> Generating missing doc.go for go packages"
	@${ROOT_DIR}/scripts/gendoc.sh

.PHONY: gen.docgo.check
gen.docgo.check: gen.docgo.doc
	@n="$$(git ls-files --others '*/doc.go' | wc -l)"; \
	if test "$$n" -gt 0; then \
		git ls-files --others '*/doc.go' | sed -e 's/^/  /'; \
		echo "$@: untracked doc.go file(s) exist in working directory" >&2 ; \
		false ; \
	fi

.PHONY: gen.docgo.add
gen.docgo.add:
	@git ls-files --others '*/doc.go' | $(XARGS) -- git add

.PHONY: gen.defaultconfigs
gen.defaultconfigs:
	@${ROOT_DIR}/scripts/gen_default_config.sh

.PHONY: gen.proto
gen.proto: $(API_DEPSRCS)

$(API_DEPSRCS): tools.verify.protoc-gen-go $(API_DEPS)
	@echo "===========> Generate protobuf files"
	@mkdir -p $(OUTPUT_DIR)
	@protoc -I $(PROTOC_INC_PATH) -I$(ROOT_DIR) \
	 --go_out=plugins=grpc:$(OUTPUT_DIR) $(@:.pb.go=.proto)
	@cp $(OUTPUT_DIR)/$(ROOT_PACKAGE)/$@ $@ || cp $(OUTPUT_DIR)/$@ $@
	@rm -rf $(OUTPUT_DIR)

.PHONY: gen.proto
 gen.proto: $(API_DEPSRCS)

.PHONY: gen.clean
gen.clean:
	@rm -rf ./api/client/{clientset,informers,listers}
	@$(FIND) -type f -name '*_generated.go' -delete
	@rm -f $(API_DEPSRCS)
