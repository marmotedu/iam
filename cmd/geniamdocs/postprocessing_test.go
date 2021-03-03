// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCleanupForInclude(t *testing.T) {
	tests := []struct {
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

func TestMarkdownPostProcessing(t *testing.T) {
	type args struct {
		cmd       *cobra.Command
		dir       string
		processor func(string) string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MarkdownPostProcessing(tt.args.cmd, tt.args.dir, tt.args.processor); (err != nil) != tt.wantErr {
				t.Errorf("MarkdownPostProcessing() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cleanupForInclude(t *testing.T) {
	type args struct {
		md string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanupForInclude(tt.args.md); got != tt.want {
				t.Errorf("cleanupForInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}
