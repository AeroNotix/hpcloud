package hpcloud

import (
	"net/http"
	"testing"
)

var createserverresponse = `
{
    "server": {
        "status": "BUILD(scheduling)",
        "updated": "2013-03-27T15:22:26Z",
        "hostId": "",
        "user_id": "22858383245118",
        "name": "MyAwesomeNewServer",
        "links": [
            {
                "href": "https://az-1.region-a.geo-1.compute.hpcloudsvc.com/v1.1/fake_tenant/servers/1032975",
                "rel": "self"
            },
            {
                "href": "https://az-1.region-a.geo-1.compute.hpcloudsvc.com/fake_tenant/servers/1032975",
                "rel": "bookmark"
            }
        ],
        "addresses": {},
        "tenant_id": "fake_tenant",
        "image": {
            "id": "1361",
            "links": [
                {
                    "href": "https://az-1.region-a.geo-1.compute.hpcloudsvc.com/fake_tenant/images/1361",
                    "rel": "bookmark"
                }
            ]
        },
        "created": "2013-03-27T15:22:26Z",
        "uuid": "c0650267-5151-4797-a2ef-cc164dfa5d6d",
        "accessIPv4": "",
        "accessIPv6": "",
        "key_name": "me",
        "adminPass": "EEB3oWJDZmMDgquj",
        "flavor": {
            "id": "100",
            "links": [
                {
                    "href": "https://az-1.region-a.geo-1.compute.hpcloudsvc.com/fake_tenant/flavors/100",
                    "rel": "bookmark"
                }
            ]
        },
        "config_drive": "",
        "id": 1032975,
        "security_groups": [
            {
                "name": "default",
                "links": [
                    {
                        "href": "https://az-1.region-a.geo-1.compute.hpcloudsvc.com/v1.1/fake_tenant/os-security-groups/73199",
                        "rel": "bookmark"
                    }
                ],
                "id": 73199
            }
        ],
        "metadata": {}
    }
}`

func TestCreateServerPrerequisites(t *testing.T) {
	httpTestsSetUp(nil)
	_, err := test_account.CreateServer(Server{
		ImageRef: DebianSqueeze6_0_3Kernel,
		Name:     "TestServer",
	})
	if err == nil {
		t.Error("Failed to account for a blank Flavour reference.")
	}
	_, err = test_account.CreateServer(Server{
		FlavorRef: XSmall,
		Name:      "TestServer",
	})
	if err == nil {
		t.Error("Failed to account for a blank Image reference.")
	}
	_, err = test_account.CreateServer(Server{
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
	test_account.CreateServer(Server{
		ImageRef:  DebianSqueeze6_0_3Kernel,
		FlavorRef: XSmall,
		Name:      "TestServer",
	})
}

func TestCreateServerUnmarshal(t *testing.T) {
	httpTestsSetUp(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(createserverresponse))
	})
	_, err := test_account.CreateServer(Server{
		ImageRef:  DebianSqueeze6_0_3Kernel,
		FlavorRef: XSmall,
		Name:      "TestServer",
	})
	if err != nil {
		t.Error(err)
	}
}
