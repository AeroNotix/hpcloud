package hpcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a Access) ActivateCDNContainer(container string) error {
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s", CDN_URL, a.TenantID, container)
	req, err := http.NewRequest("PUT", path, nil)
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
	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	return errors.New(fmt.Sprintf("Non-201 status code: %d", resp.StatusCode))
}

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

func (a Access) UpdateCDNEnabledContainerMetadata(container string, data map[string]string) error {
	client := &http.Client{}
	path := fmt.Sprintf("%s%s/%s", CDN_URL, a.TenantID, container)
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-Token", a.AuthToken())
	if len(data) > 0 {
		for key, value := range data {
			req.Header.Add(key, value)
		}
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

func (a Access) DisableCDNEnabledContainer(container string) {

}

func (a Access) DeleteCDNEnabledContainer(container string) {

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
