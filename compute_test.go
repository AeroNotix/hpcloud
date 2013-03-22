package hpcloud

import (
	"errors"
	"net/http"
	"testing"
)

type testRoundTripper struct {
	Status  bool
	Message string
}

func setUp() {
	testclient.Transport = &testtrasport
	account.Client = testclient
	testtrasport.Status = true
	account.A.Token.ID = "faketoken"
}

func (trt *testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("X-Auth-Token") == "" {
		trt.Status = false
		trt.Message = "X-Auth-Token Missing"
	}

	return nil, errors.New("Not implemented")
}

var account Access
var testclient http.Client
var testtrasport testRoundTripper

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
	setUp()
	account.baseComputeRequest("", "", nil)
	if !testtrasport.Status {
		t.Error(testtrasport.Message)
	}
}
