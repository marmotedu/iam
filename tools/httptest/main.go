// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		message := `{"status":"ok"}`
		fmt.Fprint(w, message)
	})

	addr := ":6667"
	fmt.Printf("Serving http service on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
