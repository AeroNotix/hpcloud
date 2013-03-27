package hpcloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

/*
 ListDBInstances will list all the available database instances
*/
func (a Access) ListDBInstances() (*DBInstances, error) {
	url := fmt.Sprintf("%s%s/instances", RDB_URL, a.TenantID)
	body, err := a.baseRequest(url, "GET", nil)
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

func (a Access) ListAllFlavors() (*DBFlavors, error) {
	url := fmt.Sprintf("%s%s/flavors", RDB_URL, a.TenantID)
	body, err := a.baseRequest(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	flv := &DBFlavors{}
	err = json.Unmarshal(body, flv)
	if err != nil {
		return nil, err
	}
	return flv, nil
}

/*
 CreateDBInstance creates new database instance in the HPCloud using
settings found in the DatabaseReq instance passed to this function

 This function implements the interface as described in:
 http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/create-instance.html
*/
func (a Access) CreateDBInstance(db DatabaseReq) (*NewDBInstance, error) {
	b, err := json.M
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s/instances", RDB_URL, a.TenantID)

	body, err := a.baseRequest(url, "POST",
		strings.NewReader(string(b)))
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

type Database struct {
	Name      string `json:"name"`
	FlavorRef string `json:"flavorRef"`
	Port      int    `json:"port"`
	DBType    struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"dbtype"`
}

type DBFlavors struct {
	Flavors []DBFlavor `json:"flavors"`
}

/*
 Type describing database flavor
*/
type DBFlavor struct {
	Id    int    `json:"id"`
	Links []Link `json:"links"`
	Name  string `json:"name"`
	Ram   int    `json:"ram"`
	Vcpu  int    `json:"vcpu"`
}

/*
 This type describes the JSON data which should be sent to the 
create database instance resource.
*/
type DatabaseReq struct {
	Instance Database `json:"instance"`
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

func (f DBFlavors) GetFlavorRef(fn string) string {
	for _, val := range f.Flavors {
		if val.Name == fn {
			return val.Links.HREF
		}
	}
	panic("Flavor not found")
}

func (db DatabaseReq) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`{"instance":{`)
	if db.Instance.Name == "" {
		return nil, errors.New("A name is required")
	}
	b.WriteString(fmt.Sprintf(`"name":"%s",`, db.Instance.Name))
	if db.Instance.FlavorRef == "" {
		return nil, errors.New("Flavor is required")
	}
	b.WriteString(fmt.Sprintf(`"flavorRef":"%s",`,
		db.Instance.FlavorRef))
	if db.Instance.Port == "" {
		b.WriteString(`"port":"3306",`)
	} else {
		b.WriteString(fmt.Sprintf(`"port":"%s",`, db.Instance.Port))
	}
	b.WriteString(`"dbtype":{`)
	b.WriteString(`"name":"mysql",`)
	b.WriteString(`"version":"5.5"}}}`)

	return b.Bytes(), nil
}
