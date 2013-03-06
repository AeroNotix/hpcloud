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
