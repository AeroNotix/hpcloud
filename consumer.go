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
	resp, err := http.Post(
		REGION_URL+"tokens",
		"application/json",
		strings.NewReader(string(d)),
	)
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
