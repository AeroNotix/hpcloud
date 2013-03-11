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
