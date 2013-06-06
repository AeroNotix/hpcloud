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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
  The CDN endpoints are the most "ReSTful" of all the HPCloud endpoints,
  we use the same endpoint for each container and change the verb and
  status code we use with each one.
*/
func (a Access) baseCDNRequest(method, container string, StatusCode int) error {
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s", CDN_URL, a.TenantID, container)
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Token", a.AuthToken())
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == StatusCode {
		return nil
	}
	return errors.New(fmt.Sprintf("Non-%d status code: %d", StatusCode, resp.StatusCode))

}

/*
  Activates a container for the CDN network.
*/
func (a Access) ActivateCDNContainer(container string) error {
	return a.baseCDNRequest("PUT", container, http.StatusCreated)
}

/*
  Lists available containers.

  When enabled_only == true you will only receive the containers which
  are enabled and the disabled containers will be ignored.
*/
func (a Access) ListCDNEnabledContainers(enabled_only bool) (*CDNContainers, error) {
	client := &http.Client{}
	qstring := "?format=json"
	if enabled_only {
		qstring = qstring + "&enabled_only=true"
	}
	path := fmt.Sprintf("%s%s%s", CDN_URL, a.TenantID, qstring)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", a.AuthToken())
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		c := &CDNContainers{}
		err = json.Unmarshal(b, c)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, errors.New(fmt.Sprintf("Non-200 status code: %d", resp.StatusCode))
}

/*
  Updates the metadata associated with a container.

  data will be the extra headers sent with the request, which ultimately
  end up being the metadata.
*/
func (a Access) UpdateCDNEnabledContainerMetadata(container string, data map[string]string) error {
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s", CDN_URL, a.TenantID, container)
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Token", a.AuthToken())
	for key, value := range data {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		return nil
	}
	return errors.New(fmt.Sprintf("Non-202 status code: %d", resp.StatusCode))

}

/*
  Will return the metadata associated with a single container.
*/
func (a Access) RetrieveCDNEnabledContainerMetadata(container string) (*http.Header, error) {
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s", CDN_URL, a.TenantID, container)
	req, err := http.NewRequest("HEAD", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", a.AuthToken())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return &resp.Header, nil
	}
	return nil, errors.New(fmt.Sprintf("Non-205 status code: %d", resp.StatusCode))
}

/*
  Disables a container from the CDN. This is usually preferred over
  deleting the CDN container since the the container will remain in
  the CDN and thus can be activated at a later time with no overhead.
*/
func (a Access) DisableCDNEnabledContainer(container string) error {
	return a.UpdateCDNEnabledContainerMetadata(container, map[string]string{
		"X-CDN-Enabled": "False",
	})
}

/*
  Re-enables a container.
*/
func (a Access) EnableCDNEnabledContainer(container string) error {
	return a.UpdateCDNEnabledContainerMetadata(container, map[string]string{
		"X-CDN-Enabled": "True",
	})
}

/*
  Entirely deletes a container from the CDN, note: this does not delete
  the container from the objectstore.
*/
func (a Access) DeleteCDNEnabledContainer(container string) error {
	return a.baseCDNRequest("DELETE", container, http.StatusNoContent)
}

type CDNContainers []CDNContainer
type CDNContainer struct {
	Name         string `json:"name"`
	CDNEnabled   bool   `json:"cdn_enabled"`
	TTL          int64  `json:"ttl"`
	CDNUri       string `json:"x-cdn-uri"`
	SSLCDNUri    string `json:"x-cdn-ssl-uri"`
	LogRetention bool   `json:"log_retention"`
}
