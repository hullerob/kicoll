// © 2014, Robert Hülle

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	KindleDir            string
	KindleCollectionFile string
	PreserveLastAccess   bool
)

func init() {
	flag.StringVar(&KindleDir, "dir", "/mnt/us", "where kindle is mounted")
	flag.StringVar(&KindleCollectionFile, "col", "system/collections.json", "where is collection file")
	flag.BoolVar(&PreserveLastAccess, "la", true, "preserve last access in collections.json")
}

func main() {
	flag.Parse()
	var oc map[string]Collection
	if PreserveLastAccess {
		oc = loadOldCollections(KindleDir + "/" + KindleCollectionFile)
	}
	f := FindFiles(KindleDir)
	c := BuildCollection(f, oc)
	b, err := EncodeCollections(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not encode collections: %v\n", err)
		return
	}
	file, err := os.Create(KindleDir + "/" + KindleCollectionFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "con not open collection file: %v\n", err)
		return
	}
	defer file.Close()
	_, err = file.Write(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not write to collection file: %v\n", err)
	}
}

func loadOldCollections(path string) map[string]Collection {
	data, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not read collections.json: %v\n", err)
		return nil
	}
	col, err := DecodeCollections(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not decode json file: %v\n", err)
		return nil
	}
	return col
}
