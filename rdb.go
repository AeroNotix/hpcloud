package hpcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"ioutil"
	"net/http"
)

func (a Access) ListDBInstances() (*DBInstances, error) {
	body, err := a.baseRDBRequest("instances", "GET", nil)
	if err != nil {
		return nil, err
	}
	dbs := &DBInstances{}
	err := json.Unmarshal(body, dbs)
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (a Access) baseRDBRequest(url, method string, b io.Reader) ([]byte, error) {
	path := fmt.Sprintf("%s%s/%s", RDB_URL, a.TenantID, url)
	req, err := http.NewRequest(method, path, b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", a.AuthToken())

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	case http.ServerCreated:
		//add Created Response
	case http.StatusUnathorized:
		ua := &Unathorized{}
		err = json.Unmarshal(body, ua)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ua.Message())
	case http.StatusForbidden:
		//forbidden err
	case http.StatusInternalServerError:
		//int serv err
	case http.StatusNotFound:
		nf := &NotFound{}
		err = json.Unmarshal(body, nf)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(nf.Message())
	default:
		br := &BadRequest{}
		err = json.Unmarshal(body, br)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(br.B.Message)
	}
	panic("Unreachable")
}

type DBInstance struct {
	Created string  `json:"created"`
	Flavor  Flavor_ `json:"flavor"`
	Id      string  `json:"id"`
	Links   []Link  `json:"links"`
	Name    string  `json:"name"`
	Status  string  `json:"name"`
}

type DBInstances struct {
	Instances []DBInstance `json:"instances"`
}
