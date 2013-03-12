package hpcloud

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"os"
)

/*
  HashedFile is an io.ReadWriter which reads a file into memory whilst
  giving you access to the hashed contents and only reading the file
  once
*/
type HashedFile struct {
	MD5          hash.Hash
	FileContents *bytes.Reader
	filecontents []byte
	Length       int
}

/*
  Implements io.Writer
*/
func (h *HashedFile) Write(p []byte) (int, error) {
	i, err := io.WriteString(h.MD5, string(p))
	if err != nil {
		return i, err
	}
	h.filecontents = append(h.filecontents, p...)
	if len(h.filecontents) > h.Length {
		h.Length = len(h.filecontents)
	}
	return len(p), nil
}

/*
  Implements io.Reader
*/
func (h HashedFile) Read(p []byte) (n int, err error) {
	return h.FileContents.Read(p)
}

/*
  Returns the current hash of the file.
*/
func (h HashedFile) Hash() string {
	return fmt.Sprintf("%x", h.MD5.Sum(nil))
}

/*
  Helper function to open, hash and return an io.ReadWriter of the
  file.
*/
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
