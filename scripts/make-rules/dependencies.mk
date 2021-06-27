# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for dependencies
#

.PHONY: dependencies.run
dependencies.run: dependencies.packages dependencies.tools

.PHONY: dependencies.packages
dependencies.packages:
	@$(GO) mod tidy

.PHONY: dependencies.tools
dependencies.tools: dependencies.tools.blocker dependencies.tools.critical

.PHONY: dependencies.tools.blocker
dependencies.tools.blocker: go.build.verify $(addprefix tools.verify., $(BLOCKER_TOOLS))

.PHONY: dependencies.tools.critical
dependencies.tools.critical: $(addprefix tools.verify., $(CRITICAL_TOOLS))

.PHONY: dependencies.tools.trivial
dependencies.tools.trivial: $(addprefix tools.verify., $(TRIVIAL_TOOLS))
