// © 2014, Robert Hülle

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	bufferSize = 0x4000
)

var (
	uuidPattern = regexp.MustCompile("([0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12})")
)

func bookHash(path string) string {
	switch {
	case strings.HasSuffix(path, ".mobi"):
		return mobiHash(path)
	default:
		return pathHash(path)
	}
}

func pathHash(path string) string {
	if !strings.HasPrefix(path, "/mnt/us/") {
		path = "/mnt/us/" + path
	}
	h := sha1.Sum([]byte(path))
	hs := hex.EncodeToString(h[:])
	return "*" + hs
}

func mobiHash(path string) string {
	data := make([]byte, bufferSize)
	n := readMax(data, path)
	data = data[:n]
	hash := findUUID(data)
	if len(hash) == 0 {
		return pathHash(path)
	}
	return "#" + hash + "^EBOK"
}

func findUUID(data []byte) string {
	res := uuidPattern.FindSubmatch(data)
	if len(res) < 2 {
		return ""
	}
	return string(res[1])
}

func readMax(data []byte, path string) int {
	f, err := os.Open(path)
	n := 0
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not open file: %v\n", err)
		return 0
	}
	defer f.Close()
	for {
		nn, _ := f.Read(data)
		data = data[nn:]
		n += nn
		if nn == 0 {
			break
		}
	}
	return n
}
