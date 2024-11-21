package tests

import (
	"bytes"
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {

	device_reg_url := "http://localhost:8080/OrcShack/v1/register"

	body := []byte("{\"name\": \"test\" , \"email\": \"testa@test.net\", \"password\": \"test\"   }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Device registration failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("Device registration failed, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Device registration failed, %v.", err)
	}
}

func TestLogin(t *testing.T) {

	device_reg_url := "http://localhost:8080/OrcShack/v1/login"

	body := []byte("{ \"email\": \"frans@outer.space\", \"password\": \"frans\"   }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Login failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("Login failed, %v.", err)
	}

	defer res.Body.Close()

	PrintResponse(res, t)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Login failed, %v.", err)
	}
}

func TestLoginRateLimit(t *testing.T) {

	//for i in {1..10}; do   curl -X GET http://localhost:8080/OrcShack/v1/login     -H "Content-Type: application/json"     -d '{"email":"test@example.com", "password":"securepassword"}'; done

}
