package hpcloud

import (
	"fmt"
	"net/http"
	"testing"
)

type testRoundTripper struct {
	Status  bool
	Message string
}

func (trt testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println(req)
	return nil, nil
}

var account Access
var testclient http.Client
var testtrasport testRoundTripper

func TestCreateServerPrerequisites(t *testing.T) {
	testclient.Transport = testtrasport
	account.Client = testclient

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
