package hpcloud

import (
	"net/http"
	"net/http/httptest"
)

func httpTestsSetUp(f http.HandlerFunc) {
	REGION_URL = ts.URL + "/region/"
	TOKEN_URL = REGION_URL + "tokens"
	TENANT_URL = REGION_URL + "tenants"
	OBJECT_STORE = ts.URL + "/object_store/"
	CDN_URL = ts.URL + "/cdn/"
	COMPUTE_URL = ts.URL + "/compute"
	account.A.Token.ID = "faketoken"
	if f != nil {
		functionalTest = f
	}
}

var account Access
var th = testHandler{}
var ts = httptest.NewServer(th)
var functionalTest http.HandlerFunc

type testHandler struct {
	Status         bool
	Message        string
	FunctionalTest http.HandlerFunc
}

func (th testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	functionalTest(w, req)
}
