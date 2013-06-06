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

/* Identity */
var REGION_URL = "https://region-b.geo-1.identity.hpcloudsvc.com:35357/v2.0/"
var TOKEN_URL = REGION_URL + "tokens"
var TENANT_URL = REGION_URL + "tenants"

/* Object store */
var OBJECT_STORE = "https://region-b.geo-1.objects.hpcloudsvc.com/v1.0/"

/* CDN */
var CDN_URL = "https://region-b.geo-1.cdnmgmt.hpcloudsvc.com/v1.0/"

/* Compute */
var COMPUTE_URL = "https://az-1.region-a.geo-1.compute.hpcloudsvc.com/v1.1/"

/* RDB */
var RDB_URL = "https://region-a.geo-1.dbaas-mysql.hpcloudsvc.com/v1.0/"

/* DNS */
var DNS_URL = "https://region-a.geo-1.dns.hpcloudsvc.com/v1/"
