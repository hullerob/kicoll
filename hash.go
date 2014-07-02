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

type patternID struct {
	pattern *regexp.Regexp
	prefix  string
	suffix  string
}

var patternIDs = []patternID{
	{
		regexp.MustCompile("([0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12})"),
		"#",
		"^EBOK",
	},
	{
		regexp.MustCompile("\x12(B[0-9][0-9A-Z]{5,})"),
		"#",
		"^EBOK",
	},
}

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
	hash := findPattern(data)
	if len(hash) > 0 {
		return hash
	}
	return pathHash(path)
}

func findPattern(data []byte) string {
	for _, pid := range patternIDs {
		res := pid.pattern.FindSubmatch(data)
		if len(res) > 1 {
			return pid.prefix + string(res[1]) + pid.suffix
		}
	}
	return ""
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
