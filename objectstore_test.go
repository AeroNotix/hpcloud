package hpcloud

import (
	"net/http"
	"testing"
)

func TestObjectStoreUpload(t *testing.T) {
	httpTestsSetUp(func(w http.ResponseWriter, req *http.Request) {
		if req.Body == nil {
			t.Error("Body cannot be nil")
		}
		if req.Header.Get("Etag") == "" {
			t.Error("MD5 missing from Upload request.")
		}
		if req.Header.Get("X-Auth-Token") == "" {
			t.Error("Missing auth token.")
		}
		if req.Header.Get("Content-Type") != "image/png" {
			t.Error("Content-type header missing or incorrect.")
		}
		if req.Header.Get("fake") != "value" {
			t.Error("Did not correctly set additional headers.")
		}
		w.Header().Add("Etag", req.Header.Get("Etag"))
		w.WriteHeader(http.StatusCreated)
	})
	h := &http.Header{}
	h.Add("fake", "value")
	err := test_account.ObjectStoreUpload("testfile.png", "test_container", h)
	if err != nil {
		t.Error(err)
	}
}
