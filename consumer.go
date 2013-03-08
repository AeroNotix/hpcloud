package hpcloud

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Authenticate(user, pass, tenantID string) (*Access, error) {
	l := Login{
		auth{
			credentials{
				User: user, Pass: pass,
			},
			tenantID,
		},
	}
	d, err := json.Marshal(l)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp, err := http.Post(REGION_URL+"tokens", "application/json", strings.NewReader(string(d)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	a := &Access{}

	switch resp.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, a)
		if err != nil {
			return nil, err
		}
		return a, nil
	case http.StatusBadRequest:
		b := BadRequest{}
		err = json.Unmarshal(body, &b)
		if err != nil {
			return nil, err
		}
		a.Fail = b
		return a, nil
	case http.StatusUnauthorized:
		u := Unauthorized{}
		err = json.Unmarshal(body, &u)
		if err != nil {
			return nil, err
		}
		a.Fail = u
		return a, nil
	}
	panic("Unreachable!")
}

func (a *Access) GetTenants() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", REGION_URL+"tenants", nil)
	req.Header.Add("X-Auth-Token", a.A.Token.ID)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	t := &Tenants{}
	err = json.Unmarshal(b, t)
	if err != nil {
		fmt.Println(err)
		return
	}
	a.Tenants = t.T
}

func (a Access) TenantForName(name string) (string, error) {
	for _, tenant := range a.Tenants {
		if tenant.Name == name {
			return tenant.ID, nil
		}
	}
	return "", errors.New("No tenant ID for the supplied name.")
}

func (a Access) ScopeToken(name string) (*Access, error) {
	t := TenantScope{
		Scope{
			name,
			SubToken{
				ID: a.Token(),
			},
		},
	}
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(REGION_URL+"tokens", "application/json", strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	newa := &Access{}
	switch resp.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, newa)
		if err != nil {
			return nil, err
		}
		return newa, nil
	case http.StatusBadRequest:
		b := BadRequest{}
		err = json.Unmarshal(body, &b)
		if err != nil {
			return nil, err
		}
		newa.Fail = b
		return newa, nil
	case http.StatusUnauthorized:
		u := Unauthorized{}
		err = json.Unmarshal(body, &u)
		if err != nil {
			return nil, err
		}
		newa.Fail = u
		return newa, nil
	}
	panic("Unreachable!")
}

func (a Access) ObjectStoreUpload(fpath, tenant, endpoint string) {
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fi, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	size := fi.Size()
	ssize := strconv.FormatInt(size, 10)
	fmt.Println(ssize)
	_, err = ioutil.ReadAll(f)
	fmt.Println(size)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := bytes.NewBufferString("")
	bwriter := multipart.NewWriter(buf)
	defer bwriter.Close()
	file_writer, err := bwriter.CreateFormFile("file1", fpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fh, err := os.Open(fpath)
	if err != nil {
		fmt.Println("error opening file")
		return
	}
	io.Copy(file_writer, fh)

	expires := strconv.FormatInt(time.Now().Add(1).Unix(), 10)
	path := fmt.Sprintf("/v1/%s/%s", tenant, endpoint)
	signature := a.HMAC_PostBody("104857600", "1", path, "http://google.com", expires, tenant)
	bwriter.WriteField("redirect", "http://google.com")
	bwriter.WriteField("max_file_size", "104857600")
	bwriter.WriteField("max_file_count", "10")
	bwriter.WriteField("expires", expires)
	bwriter.WriteField("signature", signature)

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf(
		"%s%s/%s", OBJECT_STORE, tenant, endpoint,
	), buf)
	req.Header.Add("X-Auth-Token", a.A.Token.ID)
	req.Header.Add("Content-type", "multipart/form-data; boundary="+bwriter.Boundary())
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	b, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(b))
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("here", err)
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

/*
 Generates the FilePOST body which should be hashed with using the
 HMAC-SHA1 hash and used as the signature for the POST request.
*/
func (a Access) HMAC_PostBody(max_file_size, max_file_count, path, redirect, expires, tenant string) string {
	bdy := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		path, redirect, max_file_size, max_file_count, expires,
	)
	return a.HMAC(a.SecretKey, tenant, bdy)
}

func (a Access) HMAC(secret_key, tenant, hmac_body string) string {
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(secret_key))
	io.WriteString(h, hmac_body)
	return fmt.Sprintf("%s:%s:%x", tenant, a.AccessKey, h.Sum(nil))
}
