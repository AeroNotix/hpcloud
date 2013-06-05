// Copyright (c) 2013, Aaron France
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.

//     * Redistributions in binary form must reproduce the above
//       copyright notice, this list of conditions and the following
//       disclaimer in the documentation and/or other materials provided
//       with the distribution.

//     * Neither the name of Aaron France nor the names of its
//       contributors may be used to endorse or promote products derived
//       from this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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
