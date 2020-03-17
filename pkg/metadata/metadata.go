// Copyright 2020 The chromium-flags Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	metadataURI        = "https://chromium.googlesource.com/chromium/src/+/refs/heads/master/chrome/browser/flag-metadata.json?format=text"
	flagDescriptionURI = "https://chromium.googlesource.com/chromium/src/+/refs/heads/master/chrome/browser/flag_descriptions.cc?format=text"
)

type Metadata struct {
	Name string `json:"name"`
}

func decodeBase64(rc io.ReadCloser) ([]byte, error) {
	r := base64.NewDecoder(base64.StdEncoding, rc)

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GetMetadata() ([]*Metadata, error) {
	resp, err := http.Get(metadataURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := decodeBase64(resp.Body)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	for scan.Scan() {
		b := scan.Bytes()
		if bytes.Contains(b, []byte("//")) {
			continue
		}
		buf.Write(b)
	}
	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading input: %v", err)
	}

	var metadatas []*Metadata
	dec := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	if err := dec.Decode(&metadatas); err != nil {
		return nil, err
	}

	return metadatas, nil
}

func GetDescription() ([]byte, error) {
	resp, err := http.Get(flagDescriptionURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := decodeBase64(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
