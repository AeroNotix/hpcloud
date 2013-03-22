package hpcloud

import (
	"net/http"
	"testing"
)

var ValidAuthenticateResponse string = `
{"access": {
   "token":    {
      "expires": "2011-10-14T21:42:59.455Z",
      "id": "faketoken",
      "tenant":       {
         "id": "14541255461800",
         "name": "HR Tenant Services"
      }
   },
   "user":    {
      "id": "30744378952176",
      "name": "arunkant",
      "roles":       [
                  {
            "id": "00000000004008",
            "serviceId": "120",
            "name": "nova:developer",
            "tenantId": "14541255461800"
         }
      ]
   },
   "serviceCatalog":    [
            {
         "name": "storage5063096349006363528",
         "type": "compute",
         "endpoints": [         {
            "adminURL": "https://localhost:8443/identity/api/v2.0/admin/0",
            "internalURL": "https://localhost:8443/identity/api/v2.0/internal/0",
            "publicURL": "https://localhost:8443/identity/api/v2.0/public/0",
            "region": "SFO"
         }]
      }
   ]
}}
`

func TestAuthenticate(t *testing.T) {
	httpTestsSetUp(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(ValidAuthenticateResponse))
	})
	acc, err := Authenticate("", "", "")
	if err != nil {
		t.Error(err)
	}
	if acc.AuthToken() != "faketoken" {
		t.Error("Failed to correctly parse the authentication response.")
	}
}

func TestAuthenticate400(t *testing.T) {
	httpTestsSetUp(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	acc, err := Authenticate("", "", "")
	if err == nil {
		t.Error("Failed to properly handle 400.")
	}
	if acc != nil {
		t.Error("Send back a useable account when the authenticate call failed.")
	}
}

func TestAuthenticate401(t *testing.T) {
	httpTestsSetUp(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	acc, err := Authenticate("", "", "")
	if err == nil {
		t.Error("Failed to properly handle 401.")
	}
	if acc != nil {
		t.Error("Send back a useable account when the authenticate call failed.")
	}
}

func TestAuthenticate403(t *testing.T) {
	httpTestsSetUp(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	acc, err := Authenticate("", "", "")
	if err == nil {
		t.Error("Failed to properly handle 403.")
	}
	if acc != nil {
		t.Error("Send back a useable account when the authenticate call failed.")
	}
}

func TestAuthenticate500(t *testing.T) {
	httpTestsSetUp(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	acc, err := Authenticate("", "", "")
	if err == nil {
		t.Error("Failed to properly handle 500.")
	}
	if acc != nil {
		t.Error("Send back a useable account when the authenticate call failed.")
	}
}
