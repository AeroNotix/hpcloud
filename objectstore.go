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
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

/*
 ObjectStoreUpload allows you to upload a file onto the HPCloud, it will
 hash the file and check the returned hash to ensure end-to-end integrity.

 It also takes an optional header which will have it's contents added
 to the request.
*/
func (a Access) ObjectStoreUpload(filename, container string, header *http.Header) error {
	f, err := OpenAndHashFile(filename)
	if err != nil {
		return err
	}
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s/%s", OBJECT_STORE, a.TenantID, container, filepath.Base(filename))
	req, err := http.NewRequest("PUT", path, f)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	req.Header.Add("Etag", f.Hash())
	req.Header.Add("X-Auth-Token", a.AuthToken())
	if header != nil {
		for key, value := range *header {
			for _, s := range value {
				req.Header.Add(key, s)
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.Header.Get("Etag") != f.Hash() {
		return errors.New("MD5 hashes do not match. Integrity not guaranteed.")
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Non-201 status code: %d", resp.StatusCode))
	}
	return nil
}

func (a Access) ObjectStoreDelete(filename string) error {
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s", OBJECT_STORE, a.TenantID, filename)
	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Token", a.AuthToken())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Non-204 status code: %d", resp.StatusCode))
	}
	return nil
}

func (a Access) ListObjects(directory string) (*FileList, error) {
	path := fmt.Sprintf("%s%s/%s", OBJECT_STORE, a.TenantID, directory)
	body, err := a.baseRequest(path, "GET", nil)
	fl := &FileList{}
	err = json.Unmarshal(body, fl)
	if err != nil {
		return nil, err
	}
	// TODO: Put in the date parsing here.
	return fl, nil
}

/*
 TemporaryURL will generate the temporary URL for the supplied filename.
*/
func (a Access) TemporaryURL(filename, expires string) string {
	hmac_path := fmt.Sprintf("/v1.0/%s/%s", a.TenantID, filename)
	hmac_body := fmt.Sprintf("%s\n%s\n%s", "GET", expires, hmac_path)
	return fmt.Sprintf("%s%s/%s?temp_url_sig=%s&temp_url_expires=%s",
		OBJECT_STORE, a.TenantID, filename, a.HMAC(a.SecretKey, a.TenantID, hmac_body),
		expires,
	)
}

type File struct {
	Hash            string `json:"hash"`
	StrLastModified string `json:"last_modified"`
	LastModified    *time.Time
	Bytes           int64  `json:"bytes"`
	Name            string `json:"name"`
	ContentType     string `json:"content_type"`
}

type FileList []File
