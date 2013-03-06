package hpcloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var REGION_URL = "https://region-b.geo-1.identity.hpcloudsvc.com:35357/v2.0/"

func Authenticate(user, pass string) *Access {
	l := Login{auth{credentials{User: user, Pass: pass}}}
	d, err := json.Marshal(l)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(d))
	resp, err := http.Post(REGION_URL+"tokens", "application/json", strings.NewReader(string(d)))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	a := &Access{}
	err = json.Unmarshal(body, a)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return a
}
