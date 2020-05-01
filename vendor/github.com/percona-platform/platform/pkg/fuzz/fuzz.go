// Package fuzz provides fuzzing helpers.
package fuzz

import (
	"crypto/sha1" //nolint:gosec
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

//nolint:gochecknoglobals
var corpusM sync.Mutex

// AddToCorpus adds data to go-fuzz corpus.
func AddToCorpus(prefix string, b []byte) {
	corpusM.Lock()
	defer corpusM.Unlock()

	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller failed")
	}
	dir := filepath.Join(filepath.Dir(file), "fuzzdata", "corpus")
	if err := os.MkdirAll(dir, 0750); err != nil {
		panic(err)
	}

	// go-fuzz uses SHA1 for non-cryptographic hashing
	file = fmt.Sprintf("%040x", sha1.Sum(b)) //nolint:gosec
	if prefix != "" {
		prefix = strings.Replace(prefix, " ", "_", -1)
		prefix = strings.Replace(prefix, "/", "_", -1)
		file = prefix + "-" + file
	}

	path := filepath.Join(dir, file)
	if err := ioutil.WriteFile(path, b, 0640); err != nil { //nolint:gosec
		panic(err)
	}
}
