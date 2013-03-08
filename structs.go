package hpcloud

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

type AuthResponse struct {
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

func (a Access) Token() string {
	return a.A.Token.ID
}

type BadRequest struct {
	B struct {
		Message string `json:"message"`
		Details string `json:"details"`
		Code    int64  `json:"code"`
	} `json:"BadRequest"`
}

type Unauthorized struct {
	U struct {
		Code            int64  `json:"code"`
		Details         string `json:"details"`
		Message         string `json:"message"`
		OtherAttributes struct {
		} `json:"otherAttributes"`
	} `json:"unauthorized"`
}

type SubToken struct {
	ID string `json:"id"`
}

type Scope struct {
	TenantName string   `json:"tenantName"`
	S          SubToken `json:"token"`
}

type TenantScope struct {
	S Scope `json:"auth"`
}

type Tenant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
}

type Tenants struct {
	T []Tenant `json:"tenants"`
}

type FailureResponse interface {
	Code() int64
	Details() string
	Message() string
}

func (u Unauthorized) Code() int64 {
	return u.U.Code
}

func (b BadRequest) Code() int64 {
	return b.B.Code
}

func (u Unauthorized) Details() string {
	return u.U.Details
}

func (b BadRequest) Details() string {
	return b.B.Details
}

func (u Unauthorized) Message() string {
	return u.U.Message
}

func (b BadRequest) Message() string {
	return b.B.Message
}
