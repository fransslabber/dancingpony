package tests

import (
	"bytes"
	"net/http"
	"testing"
)

func PrintResponse(res *http.Response, t *testing.T) {
	resbod := []byte{}
	tmpbuf := make([]byte, 256)
	offset := 0
	for {
		numred, _ := res.Body.Read(tmpbuf)
		if numred == 0 {
			break
		}
		resbod = append(resbod, tmpbuf...)
		offset += numred
	}
	t.Logf("Response Status: %v Body: %v", res.StatusCode, string(resbod[:offset]))
}

func TestSearchDish(t *testing.T) {

	device_reg_url := "http://localhost:8080/OrcShack/v1/search_dish"

	body := []byte("{ \"name\": \"Koni\"  }")
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

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Login failed, %v.", err)
	} else {
		PrintResponse(res, t)
	}

}

func TestCreateDish(t *testing.T) {

	device_reg_url := "http://localhost:8080/OrcShack/v1/create_dish"

	body := []byte("{ \"name\": \"Koni Balls Pizza\" , \"description\": \"Juicy koni balls with blue cheese.\" , \"price\": 8.99, \"category\": \"main course\", \"is_vegetarian\": false,\"is_available\": true,\"rating\": 4.98 }")

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

	if res.StatusCode != http.StatusOK {
		PrintResponse(res, t)
		t.Fatalf("Login failed, %v.", err)
	}
}
