// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// MarkdownPostProcessing goes though the generated files.
func MarkdownPostProcessing(cmd *cobra.Command, dir string, processor func(string) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		if err := MarkdownPostProcessing(c, dir, processor); err != nil {
			return err
		}
	}

	basename := strings.ReplaceAll(cmd.CommandPath(), " ", "_") + ".md"
	filename := filepath.Join(dir, basename)

	markdownBytes, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		return err
	}

	processedMarkDown := processor(string(markdownBytes))

	return ioutil.WriteFile(filename, []byte(processedMarkDown), 0o600)
}

// cleanupForInclude parts of markdown that will make difficult to use it as include in the website:
// - The title of the document (this allow more flexibility for include, e.g. include in tabs)
// - The sections see also, that assumes file will be used as a main page.
func cleanupForInclude(md string) string {
	lines := strings.Split(md, "\n")

	cleanMd := ""

	for i, line := range lines {
		if i == 0 {
			continue
		}

		if line == "### SEE ALSO" {
			break
		}

		cleanMd += line
		if i < len(lines)-1 {
			cleanMd += "\n"
		}
	}

	return cleanMd
}
