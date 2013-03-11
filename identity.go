package hpcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	a.TenantID = tenantID
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
	resp, err := http.Post(
		REGION_URL+"tokens",
		"application/json",
		strings.NewReader(string(b)),
	)
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

type Access struct {
	A struct {
		Token    Token            `json:"token"`
		User     User             `json:"user"`
		Catalogs []ServiceCatalog `json:"serviceCatalog"`
	} `json:"access"`
	Fail      FailureResponse
	Tenants   []Tenant
	SecretKey string
	AccessKey string
	TenantID  string
}

type Login struct {
	Auth auth `json:"auth"`
}

type auth struct {
	Creds    credentials `json:"passwordCredentials"`
	TenantID string      `json:"tenantId"`
}

type credentials struct {
	User string `json:"username"`
	Pass string `json:"password"`
}

type Endpoint struct {
	PublicURL   string `json:"publicURL"`
	Region      string `json:"region"`
	VersionID   string `json:"versionId"`
	VersionList string `json:"versionList"`
}

type ServiceCatalog struct {
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Role struct {
	ID        string `json:"id"`
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

type Token struct {
	Expires string      `json:"expires"`
	ID      string      `json:"id"`
	Tenant  interface{} `json:"tenant"`
}

func (a Access) Token() string {
	return a.A.Token.ID
}

type Scope struct {
	TenantName string   `json:"tenantName"`
	S          SubToken `json:"token"`
}

type TenantScope struct {
	S Scope `json:"auth"`
}

/*
 Tenant describes the response which is returned from any resource
 which contains Tenant information
*/
type Tenant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}
