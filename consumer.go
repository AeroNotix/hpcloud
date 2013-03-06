package hpcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var REGION_URL = "https://region-b.geo-1.identity.hpcloudsvc.com:35357/v2.0/"

func Authenticate(user, pass string) (*Access, FailureResponse, error) {
	l := Login{auth{credentials{User: user, Pass: pass}}}
	d, err := json.Marshal(l)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	resp, err := http.Post(REGION_URL+"tokens", "application/json", strings.NewReader(string(d)))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		a := &Access{}
		err = json.Unmarshal(body, a)
		if err != nil {
			return nil, nil, err
		}
		return a, nil, nil
	case http.StatusBadRequest:
		b := BadRequest{}
		err = json.Unmarshal(body, &b)
		if err != nil {
			return nil, nil, err
		}
		return nil, b, nil
	case http.StatusUnauthorized:
		u := Unauthorized{}
		err = json.Unmarshal(body, &u)
		if err != nil {
			return nil, nil, err
		}
		return nil, u, nil
	}
	panic("Unreachable!")
}
