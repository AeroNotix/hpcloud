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
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Domain struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	TTL       int64   `json:"ttl"`
	Serial    int64   `json:"serial"`
	Email     string  `json:"email"`
	CreateAt  string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	Data      string  `json:"data"`
}

/* Record and Domain have the same fields. */
type Record Domain

type DNSError struct {
	Path      string `json:"path"`
	Message   string `json:"message"`
	Validator string `json:"validator"`
}

type DNSErrorResponse struct {
	Message string     `json:"message"`
	Code    int        `json:"code"`
	Type    string     `json:"type"`
	Errors  []DNSError `json:"errors"`
}

func (a Access) CreateDomain(name, email string, ttl int64) (*Domain, error) {
	type BaseDNSRequest struct {
		Name  string `json:"name"`
		TTL   int64  `json:"ttl"`
		Email string `json:"email"`
	}
	d := &BaseDNSRequest{
		Name:  name,
		Email: email,
		TTL:   ttl,
	}
	jsonbody, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	body, err := a.baseDNSRequest(
		fmt.Sprintf("%sdomains", DNS_URL),
		"POST",
		strings.NewReader(string(jsonbody)),
	)
	if err != nil {
		return nil, err
	}
	d2 := &Domain{}
	err = json.Unmarshal(body, d2)
	return d2, err
}

func (a Access) DeleteDomain(domain Domain) error {
	_, err := a.baseDNSRequest(
		fmt.Sprintf("%sdomains/%s", DNS_URL, domain.ID),
		"DELETE",
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (a Access) ListDomains() ([]Domain, error) {
	type Domains struct {
		D []Domain `json:"domains"`
	}
	d := &Domains{}
	body, err := a.baseDNSRequest(
		fmt.Sprintf("%sdomains", DNS_URL),
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &d)
	return d.D, err
}

func (a Access) CreateRecord(domain Domain, t, data string) (*Record, error) {
	type DNSRecordRequest struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Data string `json:"data"`
	}
	drr := &DNSRecordRequest{
		Name: domain.Name,
		Type: t,
		Data: data,
	}
	jsonbody, err := json.Marshal(drr)
	if err != nil {
		return nil, err
	}
	body, err := a.baseDNSRequest(
		fmt.Sprintf("%sdomains/%s/records", DNS_URL, domain.ID),
		"POST",
		strings.NewReader(string(jsonbody)),
	)
	if err != nil {
		return nil, err
	}
	r := &Record{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (a Access) baseDNSRequest(url, method string, b io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}
	if a.Authenticated {
		req.Header.Add("X-Auth-Token", a.AuthToken())
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case
		http.StatusNoContent,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusOK:
		return body, nil
	case http.StatusNotFound:
		nf := &NotFound{}
		err = json.Unmarshal(body, nf)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(nf.Message())
	case
		http.StatusBadRequest,
		http.StatusConflict:
		br := &DNSErrorResponse{}
		err = json.Unmarshal(body, br)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(br.Message)
	case http.StatusUnauthorized:
		u := &Unauthorized{}
		err = json.Unmarshal(body, &u)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(u.Message())
	case http.StatusForbidden:
		f := &Forbidden{}
		err = json.Unmarshal(body, &f)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(f.Message())
	case http.StatusInternalServerError:
		ise := &InternalServerError{}
		err = json.Unmarshal(body, &ise)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ise.Message())
	default:
		panic(fmt.Sprintf("Unhandled response type: %d", resp.StatusCode))
	}
	panic("Unreachable!")
}
