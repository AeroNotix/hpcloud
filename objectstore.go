package hpcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
func (a Access) ObjectStoreUpload(filename, container, as string, header *http.Header) error {
	f, err := OpenAndHashFile(filename)
	if err != nil {
		return err
	}
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s/%s", OBJECT_STORE, a.TenantID, container, as)
	req, err := http.NewRequest("PUT", path, f)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	req.Header.Add("Etag", f.Hash())
	req.Header.Add("X-Auth-Token", a.AuthToken())
	if err != nil {
		return err
	}
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
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Non-204 status code: %d", resp.StatusCode))
	}
	return nil
}

func (a Access) ListObjects(directory string) (*FileList, error) {
	path := fmt.Sprintf("%s%s/%s", OBJECT_STORE, a.TenantID, directory)
	client := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", a.AuthToken())
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fl := &FileList{}
	err = json.Unmarshal(b, fl)
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
