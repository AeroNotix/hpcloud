package hpcloud

type Login struct {
	Auth auth `json:"auth"`
}

type auth struct {
	Creds credentials `json:"passwordCredentials"`
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
	ID      string      `json:"expires"`
	Tenant  interface{} `json:"tenant"`
}

type Access struct {
	A struct {
		Token    `json:"token"`
		User     `json:"user"`
		Catalogs []ServiceCatalog `json:"serviceCatalog"`
	} `json:"access"`
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
