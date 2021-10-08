package main

import (
	"bytes"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"io/ioutil"
	"os"
	"path/filepath"
)

var testData = []byte("0123456789ABCDEF")
var testPath = filepath.Join(os.TempDir(), "testdata")

func init() {
	//f := openFile(os.O_RDWR | os.O_CREATE | os.O_TRUNC)
	//f.Write(testData)
	//f.Close()
}

func main() {
	TestMmap()
}

func openFile(flags int) *os.File {
	f, err := os.OpenFile(testPath, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func TestMmap() {
	f := openFile(os.O_RDWR)
	defer f.Close()
	mmap, err := mmap.Map(f, mmap.RDWR, 0)
	println("addr: %p", &(mmap[0]))
	if err != nil {
		fmt.Errorf("error mapping: %s", err)
	}
	defer mmap.Unmap()
	if !bytes.Equal(testData, mmap) {
		fmt.Errorf("mmap != testData: %q, %q", mmap, testData)
	}

	mmap[9] = 'X'
	mmap.Flush()

	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Errorf("error reading file: %s", err)
	}
	if !bytes.Equal(fileData, []byte("012345678XABCDEF")) {
		fmt.Errorf("file wasn't modified")
	}

	// leave things how we found them
	mmap[9] = '9'
	mmap.Flush()
}
