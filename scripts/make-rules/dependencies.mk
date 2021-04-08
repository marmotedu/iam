# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for dependencies
#

.PHONY: dependencies.run
dependencies.run: dependencies.critical dependencies.prefer

.PHONY: dependencies.critical
dependencies.critical: go.build.verify go.lint.verify release.gsemver.verify

.PHONY: dependencies.prefer
dependencies.prefer: release.git-chglog.verify release.github-release.verify
