// Copyright (c) 2013, Bartosz Kratochwil
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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

/*
 ListDBInstances will list all the available database instances

This function implements the interface as described in:
http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/list-database-instances.html
*/
func (a Access) ListDBInstances() (*DBInstances, error) {
	url := fmt.Sprintf("%s%s/instances", RDB_URL, a.TenantID)
	body, err := a.baseRequest(url, "GET", nil)
	if err != nil {
		return nil, err
	}
	dbs := &DBInstances{}
	err = json.Unmarshal(body, dbs)
	return dbs, err
}

/*
 This function takes instance ID and deletes database instance with this ID.

This function implements the interface as described in:
http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/delete-instance.html
*/
func (a Access) DeleteDBInstance(instanceID string) error {
	url := fmt.Sprintf("%s%s/instances/%s", RDB_URL, a.TenantID,
		instanceID)
	_, err := a.baseRequest(url, "DELETE", nil)
	return err
}

/*
 This function takes instance ID and restarts DB instance with this ID.

This function implements the interface as described in:
http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/restart-instance.html
*/
func (a Access) RestartDBInstance(instanceID string) error {
	b := `{"restart":{}}`
	url := fmt.Sprintf("%s%s/instances/%s/action", RDB_URL,
		a.TenantID, instanceID)
	_, err := a.baseRequest(url, "POST", strings.NewReader(b))
	return err
}

/*
 ListAllFlavors lists all available database flavors.
 This function implements interface as described in:-
http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/list-flavors.html
*/
func (a Access) ListAllFlavors() (*DBFlavors, error) {
	url := fmt.Sprintf("%s%s/flavors", RDB_URL, a.TenantID)
	body, err := a.baseRequest(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	flv := &DBFlavors{}
	err = json.Unmarshal(body, flv)
	return flv, err
}

/*
 This function returns flavor specs for given flavor.

 This function implements the interface as described in:
 http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/get-flavor.html
*/
func (a Access) GetDBFlavor(ID string) (*DBFlavor, error) {
	url := fmt.Sprintf("%s/flavors/%s", RDB_URL, ID)
	body, err := a.baseRequest(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	flv := &DBFlavor{}
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
	b, err := db.MarshalJSON()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s/instances", RDB_URL, a.TenantID)

	body, err := a.baseRequest(url, "POST",
		strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}

	type respDB struct {
		Instance NewDBInstance `json:"instance"`
	}

	sr := &respDB{}
	err = json.Unmarshal(body, sr)
	if err != nil {
		return nil, err
	}
	return &sr.Instance, nil
}

/*
 This function retrieves details of the instance with provided ID.

This function implements the interface as described in:
http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/get-instance.html
*/
func (a Access) GetDBInstance(id string) (*InstDetails, error) {
	url := fmt.Sprintf("%s%s/instances/%s", RDB_URL, a.TenantID, id)
	body, err := a.baseRequest(url, "GET", nil)
	if err != nil {
		return nil, err
	}
	det := &InstDetails{}
	err = json.Unmarshal(body, det)
	return det, err
}

/*
 This function takes instance ID and resets password for this instance. It
 returns a new instance password.

 This function implements the interface as decribed in:
 http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/reset-instance-password.html
*/
func (a Access) ResetDBPassword(id string) (*DBCredentials, error) {
	b := `{"reset-password":{}}`
	url := fmt.Sprintf("%s%s/instances/%s/action", RDB_URL,
		a.TenantID, id)
	body, err := a.baseRequest(url, "POST", strings.NewReader(b))

	sr := &DBCredentials{}
	err = json.Unmarshal(body, sr)
	if err != nil {
		return nil, err
	}
	return sr, nil
}

/*
 This function lists all the security groups available for tenant.

 This function implements the interface as described in:
 http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/list-security-groups.html
*/
func (a Access) GetDBSecurityGroups() (*[]SecurityGroup, error) {
	url := fmt.Sprintf("%s%s/security-groups", RDB_URL, a.TenantID)
	body, err := a.baseRequest(url, "GET", nil)

	type resp struct {
		SecurityGroups []SecurityGroup `json:"security_groups"`
	}
	sr := &resp{}
	err = json.Unmarshal(body, sr)
	if err != nil {
		return nil, err
	}
	return &sr.SecurityGroups, nil
}

/*
 This function lists specific security group.

 This function implements the interface as described in:
 http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/get-security-group.html
*/
func (a Access) DBSecGroupDetails(sg string) (*SecurityGroup, error) {
	url := fmt.Sprintf("%s%s/security-groups/%s", RDB_URL, a.TenantID, sg)
	body, err := a.baseRequest(url, "GET", nil)

	type resp struct {
		SecurityGroup SecurityGroup `json:"security_group"`
	}

	sr := &resp{}
	err = json.Unmarshal(body, sr)
	if err != nil {
		return nil, err
	}
	return &sr.SecurityGroup, nil
}

/*
 Creates new security group rule

 This function implements the interface as described in:
 http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/create-security-group-rule.html
*/
func (a Access) CreateDBSecRule(Req DBSecRuleReq) (*DBSecRule, error) {
	url := fmt.Sprintf("%s%s/security-group-rules", RDB_URL, a.TenantID)
	b, err := Req.MarshalJSON()

	if err != nil {
		return nil, err
	}

	body, err := a.baseRequest(url, "POST",
		strings.NewReader(string(b)))

	type resp struct {
		SecurityGroupRule DBSecRule `json:"security_group_rule"`
	}

	sr := &resp{}
	err = json.Unmarshal(body, sr)
	if err != nil {
		return nil, err
	}

	return &sr.SecurityGroupRule, nil
}

/*
 Deletes security rule

 This function implements the interface as described in:
 http://api-docs.hpcloud.com/hpcloud-rdb-mysql/1.0/content/delete-security-group-rule.html
*/
func (a Access) RemoveDBSecRule(ruleID string) error {
	url := fmt.Sprintf("%s%s/security-group-rules/%s", RDB_URL,
		a.TenantID, ruleID)

	_, err := a.baseRequest(url, "DELETE", nil)

	if err != nil {
		return err
	}

	return nil
}

type DBInstance struct {
	Created string `json:"created"`
	Id      string `json:"id"`
	Links   []Link `json:"links"`
	Name    string `json:"name"`
	Status  string `json:"name"`
	Flavor  struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Links []Link `json:"links"`
	} `json:"flavor"`
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
 Instance Details type that is returned by server
*/
type InstDetails struct {
	Created        string          `json:"created"`
	Hostname       string          `json:"hostname"`
	ID             string          `json:"id"`
	Links          []Link          `json:"links"`
	Name           string          `json:"name"`
	Port           int             `json:"port"`
	SecurityGroups []SecurityGroup `json:"security_groups"`
	Status         string          `json:"status"`
	Updated        string          `json:"updated"`
	Flavor         struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Links []Link `json:"links"`
	} `json:"flavor"`
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
	Created    string        `json:"created"`
	Credential DBCredentials `json:"credential"`
	Flavor     struct {
		Name  string `json:"name"`
		ID    string `json:"id"`
		Links []Link `json:"links"`
	} `json:"flavor"`
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

/*
 DB Security Group Create request struct
*/
type DBSecRuleReq struct {
	SecurityGroupID string `json:"security_group_rule"`
	Cidr            string `json:"cidr"`
	FromPort        int64  `json:"from_port"`
	ToPort          int64  `json:"to_port"`
}

type DBSecRule struct {
	ID              string `json:"id"`
	SecurityGroupID string `json:"security_group_rule"`
	Cidr            string `json:"cidr"`
	FromPort        int64  `json:"from_port"`
	ToPort          int64  `json:"to_port"`
	Created         string `json:"created"`
}

func (f DBFlavors) GetFlavorRef(fn string) string {
	for _, val := range f.Flavors {
		if val.Name == fn {
			return val.Links[0].HREF
		}
	}
	panic("Flavor not found")
}

/*
 Creates JSON string for Create DB request
*/
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
	if db.Instance.Port == 0 {
		b.WriteString(`"port":"3306",`)
	} else {
		b.WriteString(fmt.Sprintf(`"port":"%s",`, db.Instance.Port))
	}
	b.WriteString(`"dbtype":{`)
	b.WriteString(`"name":"mysql",`)
	b.WriteString(`"version":"5.5"}}}`)

	return b.Bytes(), nil
}

func (rq DBSecRuleReq) MarshalJSON() ([]byte, error) {
	b := bytes.NewBufferString(`{"security_group_rule":{`)
	if rq.SecurityGroupID == "" {
		return nil, errors.New("Security group ID required")
	}
	b.WriteString(fmt.Sprintf(`"security_group_id":"%s",`,
		rq.SecurityGroupID))
	if rq.Cidr == "" {
		return nil, errors.New("Cidr is missing")
	}
	b.WriteString(fmt.Sprintf(`"cidr":"%s",`, rq.Cidr))
	if rq.FromPort == 0 {
		return nil, errors.New("from_port value is missing")
	}
	b.WriteString(fmt.Sprintf(`"from_port":%d,`, rq.FromPort))
	if rq.ToPort == 0 {
		return nil, errors.New("to_port value is missing")
	}
	b.WriteString(fmt.Sprintf(`"to_port":%d}}`, rq.ToPort))

	return b.Bytes(), nil
}
