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
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
)

/*
 Generates the FilePOST body which should be hashed with using the
 HMAC-SHA1 hash and used as the signature for the POST request.
*/
func (a Access) HMAC_PostBody(max_file_size, max_file_count, path,
	redirect, expires, tenant string) string {
	bdy := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		path, redirect, max_file_size, max_file_count, expires,
	)
	return a.HMAC(a.SecretKey, tenant, bdy)
}

/*
 HMAC is a helper method to interpolate and properly format the
 HMAC signature which is used on the HPCloud.
*/
func (a Access) HMAC(secret_key, tenant, hmac_body string) string {
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(secret_key))
	io.WriteString(h, hmac_body)
	return fmt.Sprintf("%s:%s:%x", tenant, a.AccessKey, h.Sum(nil))
}
