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
