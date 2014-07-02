// © 2014, Robert Hülle

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	KindleDir            string
	KindleCollectionFile string
)

func init() {
	flag.StringVar(&KindleDir, "dir", "/mnt/us", "where kindle is mounted")
	flag.StringVar(&KindleCollectionFile, "col", "system/collections.json", "where is collection file")
}

func main() {
	flag.Parse()
	f := FindFiles(KindleDir)
	c := BuildCollection(f)
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
