package hpcloud

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func setUp() {
	REGION_URL = ts.URL + "/region/"
	TOKEN_URL = REGION_URL + "tokens"
	TENANT_URL = REGION_URL + "tenants"
	OBJECT_STORE = ts.URL + "/object_store/"
	CDN_URL = ts.URL + "/cdn/"
	COMPUTE_URL = ts.URL + "/compute/"
	account.A.Token.ID = "faketoken"
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

func TestCreateServerPrerequisites(t *testing.T) {
	setUp()
	_, err := account.CreateServer(Server{
		ImageRef: DebianSqueeze6_0_3Kernel,
		Name:     "TestServer",
	})
	if err == nil {
		t.Error("Failed to account for a blank Flavour reference.")
	}
	_, err = account.CreateServer(Server{
		FlavorRef: XSmall,
		Name:      "TestServer",
	})
	if err == nil {
		t.Error("Failed to account for a blank Image reference.")
	}
	_, err = account.CreateServer(Server{
		ImageRef:  DebianSqueeze6_0_3Kernel,
		FlavorRef: XSmall,
	})
	if err == nil {
		t.Error("Failed to account for a missing name")
	}
}

func TestRequestHeadersAreCorrect(t *testing.T) {
	f := func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Auth-Token") == "" {
			t.Error("Missing auth token.")
		}
		if req.Header.Get("Content-type") == "" {
			t.Error("Missing content type")
		}
		if req.Header.Get("Content-type") != "application/json" {
			t.Error("Incorrect content type")
		}
	}
	functionalTest = f
	account.CreateServer(Server{
		ImageRef:  DebianSqueeze6_0_3Kernel,
		FlavorRef: XSmall,
		Name:      "TestServer",
	})
}
