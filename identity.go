package hpcloud

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
