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
	CollectionNameSuffix = "@en-US"
	CollectionNameFile   = ".collection.name"
)

var AllowedFiles = []string{".pdf", ".mobi"}

type Collection struct {
	Items      []string `json:"items"`
	LastAccess int      `json:"lastAccess"`
}

func EncodeCollections(collections map[string]Collection) ([]byte, error) {
	b, err := json.Marshal(collections)
	return b, err
}

func DecodeCollections(data []byte) (map[string]Collection, error) {
	col := make(map[string]Collection)
	err := json.Unmarshal(data, &col)
	if err != nil {
		return nil, err
	}
	return col, nil
}

func FindFiles(path string) map[string][]string {
	files := make(map[string][]string)
	col := CollectionName("", path+"/documents")
	DirWalk(col, path+"/documents", files)
	return files
}

func DirWalk(col, path string, files map[string][]string) {
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
		if col == "" {
			continue
		}
		name := path + "/" + fi.Name()
		files[name] = append(files[name], col)
	}
}

func CollectionName(old, path string) string {
	b, err := ioutil.ReadFile(path + "/" + CollectionNameFile)
	if os.IsNotExist(err) {
		return old
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not open '%s': %v\n",
			path+"/"+CollectionNameFile, err)
		return old
	}
	by := bytes.SplitN(b, []byte("\n"), 2)
	if len(by) < 1 {
		return old
	}
	return string(by[0])
}

func AllowedFile(name string) bool {
	for _, s := range AllowedFiles {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

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
		for cname, v := range collections {
			v.LastAccess = old[cname].LastAccess
		}
	}
	return collections
}
