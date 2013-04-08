package hpcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

/*
 Authenticate will send an authentication request to the HP Cloud and
 return an instance of the Access type.
*/
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
	a := &Access{}
	a.TenantID = tenantID
	body, err := a.baseRequest(
		TOKEN_URL,
		"POST",
		strings.NewReader(string(d)),
	)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, a)
	if err != nil {
		return nil, err
	}
	a.Authenticated = true
	return a, nil
}

func (a *Access) GetTenants() error {
	body, err := a.baseRequest(TENANT_URL, "GET", nil)
	t := &Tenants{}
	err = json.Unmarshal(body, t)
	if err != nil {
		return err
	}
	a.Tenants = t.T
	return nil
}

/*
 On the HP Cloud, tenants have name strings associated with them, you
 can find the tenantID associated with a name with this function
*/
func (a Access) TenantForName(name string) (string, error) {
	for _, tenant := range a.Tenants {
		if tenant.Name == name {
			return tenant.ID, nil
		}
	}
	return "", errors.New("No tenant ID for the supplied name.")
}

/*
 ScopeToken will scope or rescope an Auth Token to a different
 tenantID
*/
func (a Access) ScopeToken(name string) (*Access, error) {
	t := TenantScope{
		Scope{
			name,
			SubToken{
				ID: a.AuthToken(),
			},
		},
	}
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	body, err := a.baseRequest(
		TOKEN_URL,
		"POST",
		strings.NewReader(string(b)),
	)
	if err != nil {
		return nil, err
	}
	newa := &Access{}
	err = json.Unmarshal(body, newa)
	return newa, err
}

/*
 This function takes service name and region as parameters and returns
public URL for endpoint, that can be queried later on.
*/
func (a Access) GetEndpointURL(servName string, region string) string {

	for _, service := range a.A.Catalogs {
		if service.Name == servName {
			for _, endpoint := range service.Endpoints {
				if endpoint.Region == region {
					return endpoint.PublicURL
				}
			}
		}
	}
	panic("Service not found in this region")
}

/*
 Access describes the reponse received from the /tokens endpoint when
 posting with username and password.
*/
type Access struct {
	A struct {
		Token    Token            `json:"token"`
		User     User             `json:"user"`
		Catalogs []ServiceCatalog `json:"serviceCatalog"`
	} `json:"access"`
	Authenticated bool
	Tenants       []Tenant
	SecretKey     string
	AccessKey     string
	TenantID      string
	Client        http.Client
}

/*
 Token is a helper method to traverse the Access type to retrieve the
 auth_token
*/
func (a Access) AuthToken() string {
	return a.A.Token.ID
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

/*
 Role describes in-part the response you will receive when making
 an authentication request.

 Roles are services for with the user making an authentication request
 is authenticated to use.

 A personality that a user assumes when performing a specific set of
 operations. A role includes a set of rights and privileges.

 {
   "id": "00000000004003",
   "serviceId": "100",
   "name": "domainadmin"
 }

*/
type Role struct {
	ID        string `json:"id"`
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
}

/*

 User describes in-part the response you will receive when making
 an authentication request.

"user": {
    "id": "<tenant_id>",
    "name": "<username>",
    "roles": [<array_of_roles]
}
*/
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

/*

 Token describes in-part the response you will receive when making
 an authentication request.

 If you didn't supply a tenantID (currently this library does not
 support unscoped authorization requests.) then the tenant section
 will be null, hence using a pointer type for this field.

 "token": {
    "expires": "<token_expiry_date>",
    "id": "<your_auth_token>",
    "tenant": {
      "id": "<tenant_id>",
      "name": "<tenant_name>"
    }
*/
type Token struct {
	Expires string  `json:"expires"`
	ID      string  `json:"id"`
	Tenant  *Tenant `json:"tenant"`
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
