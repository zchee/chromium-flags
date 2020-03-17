// Copyright 2020 The chromium-flags Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/zchee/chromium-flags/pkg/metadata"
	color "github.com/zchee/color/v2"
	"github.com/zchee/strcase"
)

var _ = color.Red

func main() {
	log.SetFlags(log.Lshortfile)

	metadatas, err := metadata.GetMetadata()
	if err != nil {
		log.Fatal(err)
	}

	ss := make([]string, len(metadatas))
	for i, md := range metadatas {
		ss[i] = strcase.ToCamelCase(md.Name)
	}

	descriptions, err := metadata.GetDescription()
	if err != nil {
		log.Fatal(err)
	}

	scan := bufio.NewScanner(bytes.NewReader(descriptions))
	scan.Split(bufio.ScanLines)
parent:
	for scan.Scan() {
		text := scan.Text()
		for _, s := range ss {
			if strings.Contains(text, "k"+s+"Name[]") {
				fmt.Printf("%s: ", color.HiCyanString(strcase.ToSnakeCase(s)))
				continue parent
			}
		}

		if !strings.Contains(text, "[]") {
			text = strings.TrimSpace(text)
			switch {
			case strings.Contains(text, ";"):
				fmt.Printf("%s\n\n", strings.TrimSuffix(text, ";"))
			case strings.Contains(text, `"`):
				text, _ = strconv.Unquote(text)
				fmt.Printf("%s", text)
			default:
				fmt.Printf("%s", text)
			}
		}
	}
}
