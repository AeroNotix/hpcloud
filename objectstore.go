package hpcloud

import (
	"errors"
	"fmt"
	"net/http"
)

/*
 ObjectStoreUpload allows you to upload a file onto the HPCloud, it will
 hash the file and check the returned hash to ensure end-to-end integrity.

 It also takes an optional header which will have it's contents added
 to the request.
*/
func (a Access) ObjectStoreUpload(filename, container, as string, header *http.Header) error {
	f, err := OpenAndHashFile(filename)
	if err != nil {
		return err
	}
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s/%s", OBJECT_STORE, a.TenantID, container, as)
	req, err := http.NewRequest("PUT", path, f)
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Token", a.AuthToken())
	if header != nil {
		for key, value := range *header {
			for _, s := range value {
				req.Header.Add(key, s)
			}
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return errors.New("Not created.")
	}
	return nil
}

func (a Access) TemporaryURL(filename, expires string) string {
	hmac_path := fmt.Sprintf("/v1.0/%s/%s", a.TenantID, filename)
	hmac_body := fmt.Sprintf("%s\n%s\n%s", "GET", expires, hmac_path)
	return fmt.Sprintf("%s%s/%s?temp_url_sig=%s&temp_url_expires=%s",
		OBJECT_STORE, a.TenantID, filename, a.HMAC(a.SecretKey, a.TenantID, hmac_body),
		expires,
	)
}
