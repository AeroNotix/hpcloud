package hpcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

/*
 ListDBInstances will list all the available database instances
*/
func (a Access) ListDBInstances() (*DBInstances, error) {
	body, err := a.baseRDBRequest("instances", "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	dbs := &DBInstances{}
	err = json.Unmarshal(body, dbs)
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

/*
 CreateDBInstance creates new database instance in the HPCloud using
settings found in the DatabaseReq instance passed to this function

 This function implements the interface as described in:
 http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/create-instance.html
*/
func (a Access) CreateDBInstance(db DatabaseReq) (*NewDBInstance, error) {
	b, err := db.MarshalDBJSON()
	if err != nil {
		return nil, err
	}

	body, err := a.baseRDBRequest("instances", "POST",
		strings.NewReader(string(b)), 119)
	if err != nil {
		return nil, err
	}

	sr := &NewDBInstance{}
	err = json.Unmarshal(body, sr)
	if err != nil {
		return nil, err
	}
	return sr, nil
}

func (a Access) baseRDBRequest(url, method string, b io.Reader, conLen int) ([]byte, error) {
	path := fmt.Sprintf("%s%s/%s", RDB_URL, a.TenantID, url)
	req, err := http.NewRequest(method, path, b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", a.AuthToken())
	if conLen != nil {
		req.Header.Add("Content-Length", conLen)
	}

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
	case http.StatusCreated:
		return body, nil
	case http.StatusUnathorized:
		ua := &Unathorized{}
		err = json.Unmarshal(body, ua)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ua.Message())
	case http.StatusForbidden:
		fr := &Forbidden{}
		err = json.Unmarshal(body, fr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(fr.Message())
	case http.StatusInternalServerError:
		ise := &InternalServerError{}
		err = json.Unmarshal(body, ise)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ise.Message())
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

/*
 This type describes the JSON data which should be sent to the 
create database instance resource.
*/
type DatabaseReq struct {
	Instance Database `json:"instance"`
}

type Database struct {
	Name      string       `json:"name"`
	FlavorRef string       `json:"flavorRef"`
	Port      int          `json:"port"`
	Dbtype    DatabaseType `json:"port"`
}

type DatabaseType struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

/*
 This type describes JSON response from a successful CreateDBInstance
 call.
*/
type NewDBInstance struct {
	Created        string        `json:"created"`
	Credential     DBCredentials `json:"credential"`
	Flavor         Flavor_       `json:"flavor"`
	Hostname       string        `json:"hostname"`
	Id             string        `json:"id"`
	Links          []Link        `json:"links"`
	Name           string        `json:"name"`
	SecurityGroups []DBSecGroups `json:"security_groups"`
	Status         string        `json:"status"`
}

/*
 This type describes Database Security groups 
*/
type DBSecGroups struct {
	Id    string `json:"id"`
	Links []Link `json:"links"`
}

/*
 This type describes Database Credentials
*/
type DBCredentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
