// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestCleanupForInclude(t *testing.T) {
	var tests = []struct {
		markdown, expectedMarkdown string
	}{
		{ // first line is removed
			// Nb. first line is the title of the document, and by removing it you get
			//     more flexibility for include, e.g. include in tabs
			markdown: "line 1\n" +
				"line 2\n" +
				"line 3",
			expectedMarkdown: "line 2\n" +
				"line 3",
		},
		{ // everything after ###SEE ALSO is removed
			// Nb.  see also, that assumes file will be used as a main page (does not apply to includes)
			markdown: "line 1\n" +
				"line 2\n" +
				"### SEE ALSO\n" +
				"line 3",
			expectedMarkdown: "line 2\n",
		},
	}

	for _, rt := range tests {
		actual := cleanupForInclude(rt.markdown)
		if actual != rt.expectedMarkdown {
			t.Errorf(
				"failed cleanupForInclude:\n\texpected: %s\n\t  actual: %s",
				rt.expectedMarkdown,
				actual,
			)
		}
	}
}
