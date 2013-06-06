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
	"net/http/httptest"
)

func httpTestsSetUp(f http.HandlerFunc) {
	REGION_URL = ts.URL + "/region/"
	TOKEN_URL = REGION_URL + "tokens"
	TENANT_URL = REGION_URL + "tenants"
	OBJECT_STORE = ts.URL + "/object_store/"
	CDN_URL = ts.URL + "/cdn/"
	COMPUTE_URL = ts.URL + "/compute"
	RDB_URL = ts.URL + "/rdb"
	test_account.A.Token.ID = "faketoken"
	test_account.Authenticated = true
	if f != nil {
		functionalTest = f
	}
}

var test_account Access
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
