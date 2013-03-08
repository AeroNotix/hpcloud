package hpcloud

import (
	"bytes"
	"crypto/md5"
	"hash"
	"io"
	"os"
)

type HashedFile struct {
	MD5          hash.Hash
	FileContents *bytes.Reader
	filecontents []byte
	Length       int
}

func (h HashedFile) Write(p []byte) (int, error) {
	i, err := io.WriteString(h.MD5, string(p))
	if err != nil {
		return i, err
	}
	h.filecontents = append(h.filecontents, p...)
	h.Length = len(h.filecontents)
	return len(p), nil
}

func (h HashedFile) Read(p []byte) (n int, err error) {
	return h.FileContents.Read(p)
}

func OpenAndHashFile(filename string) (*HashedFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	hf := &HashedFile{MD5: md5.New()}
	io.Copy(hf, f)
	hf.FileContents = bytes.NewReader(hf.filecontents)
	return hf, nil
}
