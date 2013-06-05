// Copyright (c) 2013, Aaron France
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.

//     * Redistributions in binary form must reproduce the above
//       copyright notice, this list of conditions and the following
//       disclaimer in the documentation and/or other materials provided
//       with the distribution.

//     * Neither the name of Aaron France nor the names of its
//       contributors may be used to endorse or promote products derived
//       from this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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
