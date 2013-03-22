package hpcloud

import (
	"net/http"
	"testing"
)

func TestCreateServerPrerequisites(t *testing.T) {
	httpTestsSetUp(nil)
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
	httpTestsSetUp(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Auth-Token") == "" {
			t.Error("Missing auth token.")
		}
		if req.Header.Get("Content-type") == "" {
			t.Error("Missing content type")
		}
		if req.Header.Get("Content-type") != "application/json" {
			t.Error("Incorrect content type")
		}
	})
	account.CreateServer(Server{
		ImageRef:  DebianSqueeze6_0_3Kernel,
		FlavorRef: XSmall,
		Name:      "TestServer",
	})
}
