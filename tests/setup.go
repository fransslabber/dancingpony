package tests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	rest_api "biz.orcshack/menu/web"
)

type TestAuth struct {
	Session_jwt string
	User_id     uint32
	Srvr_url    string
	Email       string
	Password    string
	Dish_id     uint32
}

func Register(auth *TestAuth) error {

	device_reg_url := auth.Srvr_url + "/OrcShack/v1/register"

	body := []byte("{\"name\": \"Noname\" , \"email\": \"" + auth.Email + "\", \"password\": \"" + auth.Password + "\" }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("registration http request failed. %v", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
	res, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("registration do http request failed. %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		post := &rest_api.Error_Response{}
		json.NewDecoder(res.Body).Decode(post)
		return fmt.Errorf("Response status code = %v. %v", res.StatusCode, post)
	}
	return nil
}

func LoginUser(auth *TestAuth) error {

	device_reg_url := auth.Srvr_url + "/OrcShack/v1/login"
	body := []byte("{ \"email\": \"" + auth.Email + "\", \"password\": \"" + auth.Password + "\" }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("login(%s,f%s), %v http request failed.", auth.Email, auth.Password, err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
	res, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("login(%s,f%s) = %+v, %v http request failed.", auth.Email, auth.Password, res, err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		post := &rest_api.Error_Response{}
		json.NewDecoder(res.Body).Decode(post)
		return fmt.Errorf("Response status code = %v. %v", res.StatusCode, post)
	} else {
		post := &rest_api.LoginUser_ResponseUser{}
		derr := json.NewDecoder(res.Body).Decode(post)
		if derr != nil {
			return fmt.Errorf("Response status code = %v. JSON decode failed = %v.", res.StatusCode, derr)
		}
		auth.Session_jwt = post.JWT
	}

	return nil
}

func SetupLoginTest(tb testing.TB) (func(tb testing.TB), *TestAuth) {
	testAuth := TestAuth{Email: "frans@outer.space", Password: "frans", Srvr_url: "https://localhost:4443"}

	// if err := Register(&testAuth); err != nil {
	// 	tb.Fatalf(err.Error())
	// }
	if err := LoginUser(&testAuth); err != nil {
		tb.Fatalf(err.Error())
	}

	return func(tb testing.TB) {
	}, &testAuth
}
