// © 2014, Robert Hülle

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// CollectionNameSuffix is suffix of collection name in collections.json
	CollectionNameSuffix = "@en-US"

	// CollectionNameFile is name of file where collections are named.
	CollectionNameFile = ".collection.name"
)

// Suffixes of files considered for collections.
var AllowedFiles = []string{".pdf", ".mobi"}

// Collection is internal representation of collection.
type Collection struct {
	Items      []string `json:"items"`
	LastAccess int64    `json:"lastAccess"`
}

// EncodeCollections converts internal representation of collections
// to JSON format used by Kindle.
func EncodeCollections(collections map[string]Collection) ([]byte, error) {
	b, err := json.Marshal(collections)
	return b, err
}

// DecodeCollections is opposite of EncodeCollections.
func DecodeCollections(data []byte) (map[string]Collection, error) {
	col := make(map[string]Collection)
	err := json.Unmarshal(data, &col)
	if err != nil {
		return nil, err
	}
	return col, nil
}

// FindFiles returns map of filepath->collections
func FindFiles(path string) map[string][]string {
	files := make(map[string][]string)
	col := CollectionName(nil, path+"/documents")
	DirWalk(col, path+"/documents", files)
	return files
}

// DirWalk walks directory tree and associates files with collections.
func DirWalk(col []string, path string, files map[string][]string) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading directory '%s': %v\n", path, err)
	}
	for _, fi := range fis {
		if fi.Name()[0] == '.' {
			continue
		}
		if fi.IsDir() {
			npath := path + "/" + fi.Name()
			ncol := CollectionName(col, npath)
			DirWalk(ncol, npath, files)
			continue
		}
		if !AllowedFile(fi.Name()) {
			continue
		}
		if len(col) == 0 {
			continue
		}
		name := path + "/" + fi.Name()
		files[name] = append(files[name], col...)
	}
}

// CollectionName reads collection names for given directory and replaces
// collection names from parent directory if it finds any.
func CollectionName(old []string, path string) []string {
	b, err := ioutil.ReadFile(path + "/" + CollectionNameFile)
	if os.IsNotExist(err) {
		return old
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not open '%s': %v\n",
			path+"/"+CollectionNameFile, err)
		return old
	}
	by := bytes.Split(b, []byte("\n"))
	ncols := make([]string, 0, len(by))
	for _, line := range by {
		if len(line) > 0 {
			ncols = append(ncols, string(line))
		}
	}
	return ncols
}

// AllowedFile checks if file is to be considered for addition to collection.
func AllowedFile(name string) bool {
	for _, s := range AllowedFiles {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

// BuildCollection converts filepaths to hash/ID and adds them to collection.
// It also scans existing collections to preserve lastAccess.
func BuildCollection(files map[string][]string,
	old map[string]Collection) map[string]Collection {
	collections := make(map[string]Collection)
	hashCache := make(map[string]string)
	for file, colls := range files {
		for _, collectionName := range colls {
			cname := collectionName + CollectionNameSuffix
			hash := hashCache[file]
			if len(hash) == 0 {
				hash = bookHash(file)
				hashCache[file] = hash
			}
			items := collections[cname].Items
			items = append(items, hash)
			collections[cname] = Collection{Items: items}
		}
	}
	if old != nil {
		for cname, v := range old {
			if _, ok := collections[cname]; ok {
				co := collections[cname]
				co.LastAccess = v.LastAccess
				collections[cname] = co
			}
		}
	}
	return collections
}
