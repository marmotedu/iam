# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for release
#
#

.PHONY: release.run
release.run: release.verify release.ensure-tag
	@scripts/release.sh

.PHONY: release.verify
release.verify: tools.verify.git-chglog tools.verify.github-release tools.verify.coscmd

.PHONY: release.tag
release.tag: tools.verify.gsemver release.ensure-tag
	@git push origin `git describe --tags --abbrev=0`

.PHONY: release.ensure-tag
release.ensure-tag: tools.verify.gsemver
	@scripts/ensure_tag.sh
