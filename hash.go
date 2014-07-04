// © 2014, Robert Hülle

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	// ID in .mobi should be close to beginning of file.
	// This is length of searched data.
	bufferSize = 0x4000
)

type patternID struct {
	pattern *regexp.Regexp
	prefix  string
	suffix  string
}

// ID patterns for searching in .mobi files
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

// bookHash returns hash/ID by which kindle references books
func bookHash(path string) string {
	switch {
	case strings.HasSuffix(path, ".mobi"):
		return mobiHash(path)
	default:
		return pathHash(path)
	}
}

// pathHash computes sha1 hash of filepath.
func pathHash(path string) string {
	if !strings.HasPrefix(path, "/mnt/us/") {
		// strip non-standard prefix from path and replace it with /mnt/us
		i := len(KindleDir) + 1
		path = "/mnt/us/" + path[i:]
	}
	h := sha1.Sum([]byte(path))
	hs := hex.EncodeToString(h[:])
	return "*" + hs
}

// mobiHash returns hash/ID for .mobi files.
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

// findPattern searches for ID pattern in block of data.
func findPattern(data []byte) string {
	for _, pid := range patternIDs {
		res := pid.pattern.FindSubmatch(data)
		if len(res) > 1 {
			return pid.prefix + string(res[1]) + pid.suffix
		}
	}
	return ""
}

// readMax reads maximum available data from file, up to len(data).
func readMax(data []byte, path string) int {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not open file: %v\n", err)
		return 0
	}
	defer f.Close()
	n, err := io.ReadFull(f, data)
	if err != nil && err != io.ErrUnexpectedEOF {
		fmt.Fprintf(os.Stderr, "error reading file '%s': %v\n", path, err)
	}
	return n
}
