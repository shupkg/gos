package embed

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"sync"
)

type File struct {
	Path        string
	FileName    string
	FileSize    int64
	FileModTime int64

	Data string

	v []byte
	o sync.Once
}

func (f *File) Bytes() []byte {
	f.o.Do(func() {
		gr, _ := gzip.NewReader(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.Data)))
		f.v, _ = ioutil.ReadAll(gr)
	})
	return f.v
}
