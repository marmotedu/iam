// Copyrighaddlicense: t 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This program ensures source code files have copyright license headers.
// See usage with "addlicense -h".
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

const helpText = `Usage: addlicense [flags] pattern [pattern ...]

The program ensures source code files have copyright license headers
by scanning directory patterns recursively.

It modifies all source files in place and avoids adding a license header
to any file that already has one.

The pattern argument can be provided multiple times, and may also refer
to single files.

Flags:
`

const tmplApache = `Copyright {{.Year}} {{.Holder}}

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`

const tmplBSD = `Copyright (c) {{.Year}} {{.Holder}} All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.`

const tmplMIT = `Copyright (c) {{.Year}} {{.Holder}}

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`

const tmplMPL = `This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.`

type copyrightData struct {
	Year   string
	Holder string
}

var (
	holder    = pflag.StringP("holder", "c", "Google LLC", "copyright holder")
	license   = pflag.StringP("license", "l", "apache", "license type: apache, bsd, mit, mpl")
	licensef  = pflag.StringP("licensef", "f", "", "license file")
	year      = pflag.StringP("year", "y", fmt.Sprint(time.Now().Year()), "copyright year(s)")
	verbose   = pflag.BoolP("verbose", "v", false, "verbose mode: print the name of the files that are modified")
	checkonly = pflag.BoolP(
		"check",
		"",
		false,
		"check only mode: verify presence of license headers and exit with non-zero code if missing",
	)
	skipDirs  = pflag.StringSliceP("skip-dirs", "", nil, "regexps of directories to skip")
	skipFiles = pflag.StringSliceP("skip-files", "", nil, "regexps of files to skip")
	help      = pflag.BoolP("help", "h", false, "show this help message")
)

var patterns = struct {
	dirs  []*regexp.Regexp
	files []*regexp.Regexp
}{}

var (
	licenseTemplate = make(map[string]*template.Template)
	usage           = func() {
		fmt.Println(helpText)
		pflag.PrintDefaults()
	}
)

// nolint: gocognit // no lint
func main() {
	pflag.Usage = usage
	pflag.Parse()

	if *help {
		pflag.Usage()
		os.Exit(1)
	}

	if pflag.NArg() == 0 {
		pflag.Usage()
		os.Exit(1)
	}

	if len(*skipDirs) != 0 {
		ps, err := getPatterns(*skipDirs)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		patterns.dirs = ps
	}

	if len(*skipFiles) != 0 {
		ps, err := getPatterns(*skipFiles)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		patterns.files = ps
	}

	data := &copyrightData{
		Year:   *year,
		Holder: *holder,
	}

	var t *template.Template
	if *licensef != "" {
		d, err := ioutil.ReadFile(*licensef)
		if err != nil {
			fmt.Printf("license file: %v\n", err)
			os.Exit(1)
		}
		t, err = template.New("").Parse(string(d))
		if err != nil {
			fmt.Printf("license file: %v\n", err)
			os.Exit(1)
		}
	} else {
		t = licenseTemplate[*license]
		if t == nil {
			fmt.Printf("unknown license: %s\n", *license)
			os.Exit(1)
		}
	}

	// process at most 1000 files in parallel
	ch := make(chan *file, 1000)
	done := make(chan struct{})
	go func() {
		var wg errgroup.Group
		for f := range ch {
			f := f // https://golang.org/doc/faq#closures_and_goroutines
			wg.Go(func() error {
				// nolint: nestif
				if *checkonly {
					// Check if file extension is known
					lic, err := licenseHeader(f.path, t, data)
					if err != nil {
						fmt.Printf("%s: %v\n", f.path, err)
						return err
					}
					if lic == nil { // Unknown fileExtension
						return nil
					}
					// Check if file has a license
					isMissingLicenseHeader, err := fileHasLicense(f.path)
					if err != nil {
						fmt.Printf("%s: %v\n", f.path, err)
						return err
					}
					if isMissingLicenseHeader {
						fmt.Printf("%s\n", f.path)
						return errors.New("missing license header")
					}
				} else {
					modified, err := addLicense(f.path, f.mode, t, data)
					if err != nil {
						fmt.Printf("%s: %v\n", f.path, err)
						return err
					}
					if *verbose && modified {
						fmt.Printf("%s added license\n", f.path)
					}
				}
				return nil
			})
		}
		err := wg.Wait()
		close(done)
		if err != nil {
			os.Exit(1)
		}
	}()

	for _, d := range pflag.Args() {
		walk(ch, d)
	}
	close(ch)
	<-done
}

type file struct {
	path string
	mode os.FileMode
}

func getPatterns(patterns []string) ([]*regexp.Regexp, error) {
	patternsRe := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		patternRe, err := regexp.Compile(p)
		if err != nil {
			fmt.Printf("can't compile regexp %q\n", p)
			return nil, err
		}
		patternsRe = append(patternsRe, patternRe)
	}

	return patternsRe, nil
}

func walk(ch chan<- *file, start string) {
	_ = filepath.Walk(start, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("%s error: %v\n", path, err)
			return nil
		}
		if fi.IsDir() {
			for _, pattern := range patterns.dirs {
				if pattern.MatchString(fi.Name()) {
					return filepath.SkipDir
				}
			}

			return nil
		}

		for _, pattern := range patterns.files {
			if pattern.MatchString(fi.Name()) {
				return nil
			}
		}

		ch <- &file{path, fi.Mode()}
		return nil
	})
}

func addLicense(path string, fmode os.FileMode, tmpl *template.Template, data *copyrightData) (bool, error) {
	var lic []byte
	var err error
	lic, err = licenseHeader(path, tmpl, data)
	if err != nil || lic == nil {
		return false, err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil || hasLicense(b) {
		return false, err
	}

	line := hashBang(b)
	if len(line) > 0 {
		b = b[len(line):]
		if line[len(line)-1] != '\n' {
			line = append(line, '\n')
		}
		line = append(line, '\n')
		lic = append(line, lic...)
	}
	b = append(lic, b...)
	return true, ioutil.WriteFile(path, b, fmode)
}

// fileHasLicense reports whether the file at path contains a license header.
func fileHasLicense(path string) (bool, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil || hasLicense(b) {
		return false, err
	}
	return true, nil
}

func licenseHeader(path string, tmpl *template.Template, data *copyrightData) ([]byte, error) {
	var lic []byte
	var err error
	switch fileExtension(path) {
	default:
		return nil, nil
	case ".c", ".h":
		lic, err = prefix(tmpl, data, "/*", " * ", " */")
	case ".js", ".mjs", ".cjs", ".jsx", ".tsx", ".css", ".tf", ".ts":
		lic, err = prefix(tmpl, data, "/**", " * ", " */")
	case ".cc",
		".cpp",
		".cs",
		".go",
		".hh",
		".hpp",
		".java",
		".m",
		".mm",
		".proto",
		".rs",
		".scala",
		".swift",
		".dart",
		".groovy",
		".kt",
		".kts":
		lic, err = prefix(tmpl, data, "", "// ", "")
	case ".py", ".sh", ".yaml", ".yml", ".dockerfile", "dockerfile", ".rb", "gemfile":
		lic, err = prefix(tmpl, data, "", "# ", "")
	case ".el", ".lisp":
		lic, err = prefix(tmpl, data, "", ";; ", "")
	case ".erl":
		lic, err = prefix(tmpl, data, "", "% ", "")
	case ".hs", ".sql":
		lic, err = prefix(tmpl, data, "", "-- ", "")
	case ".html", ".xml", ".vue":
		lic, err = prefix(tmpl, data, "<!--", " ", "-->")
	case ".php":
		lic, err = prefix(tmpl, data, "", "// ", "")
	case ".ml", ".mli", ".mll", ".mly":
		lic, err = prefix(tmpl, data, "(**", "   ", "*)")
	}
	return lic, err
}

func fileExtension(name string) string {
	if v := filepath.Ext(name); v != "" {
		return strings.ToLower(v)
	}
	return strings.ToLower(filepath.Base(name))
}

var head = []string{
	"#!",                       // shell script
	"<?xml",                    // XML declaratioon
	"<!doctype",                // HTML doctype
	"# encoding:",              // Ruby encoding
	"# frozen_string_literal:", // Ruby interpreter instruction
	"<?php",                    // PHP opening tag
}

func hashBang(b []byte) []byte {
	line := make([]byte, 0, len(b))
	for _, c := range b {
		line = append(line, c)
		if c == '\n' {
			break
		}
	}
	first := strings.ToLower(string(line))
	for _, h := range head {
		if strings.HasPrefix(first, h) {
			return line
		}
	}
	return nil
}

func hasLicense(b []byte) bool {
	n := 1000
	if len(b) < 1000 {
		n = len(b)
	}
	return bytes.Contains(bytes.ToLower(b[:n]), []byte("copyright")) ||
		bytes.Contains(bytes.ToLower(b[:n]), []byte("mozilla public"))
}

// prefix will execute a license template t with data d
// and prefix the result with top, middle and bottom.
func prefix(t *template.Template, d *copyrightData, top, mid, bot string) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, d); err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if top != "" {
		fmt.Fprintln(&out, top)
	}
	s := bufio.NewScanner(&buf)
	for s.Scan() {
		fmt.Fprintln(&out, strings.TrimRightFunc(mid+s.Text(), unicode.IsSpace))
	}
	if bot != "" {
		fmt.Fprintln(&out, bot)
	}
	fmt.Fprintln(&out)
	return out.Bytes(), nil
}

// nolint: gochecknoinits
func init() {
	licenseTemplate["apache"] = template.Must(template.New("").Parse(tmplApache))
	licenseTemplate["mit"] = template.Must(template.New("").Parse(tmplMIT))
	licenseTemplate["bsd"] = template.Must(template.New("").Parse(tmplBSD))
	licenseTemplate["mpl"] = template.Must(template.New("").Parse(tmplMPL))
}
