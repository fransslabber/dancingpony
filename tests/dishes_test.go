package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type Dish struct {
	Id            uint32    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float32   `json:"price"`
	Category      string    `json:"category"`
	Is_vegetarian bool      `json:"is_vegetarian"`
	Is_available  bool      `json:"is_available"`
	Rating        float32   `json:"rating"`
	Restaurant_id uint32    `json:"restaurant_id"`
	Date_created  time.Time `json:"date_created"`
	Date_updated  time.Time `json:"date_updated"`
}

var dish Dish

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
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := "http://localhost:8080/OrcShack/v1/search_dish"

	body := []byte("{ \"name\": \"Koni\"  }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Search dishes failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("Search dishes, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Search dishes, %v.", err)
	} else {
		PrintResponse(res, t)
	}

}

func TestListDish(t *testing.T) {
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := testAuth.Srvr_url + "/OrcShack/v1/list_dishes"

	body := []byte("{ \"category\": \"main course\"  }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("List dishes failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("List dishes failed, %v.", err)
	}

	defer res.Body.Close()

	PrintResponse(res, t)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("List dishes failed, %v.", err)
	}
}

func TestCreateDish(t *testing.T) {
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := "http://localhost:8080/OrcShack/v1/create_dish"

	body := []byte("{ \"name\": \"Koni Balls Pizza\" , \"description\": \"Juicy koni balls with blue cheese.\" , \"price\": 8.99, \"category\": \"main course\", \"is_vegetarian\": false,\"is_available\": true,\"rating\": 4.98 }")

	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Create dish failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("Create dish failed, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		PrintResponse(res, t)
		t.Fatalf("Create dish failed, %v.", err)
	} else {
		//PrintResponse(res, t)
		err = json.NewDecoder(res.Body).Decode(&dish)
		if err != nil {
			t.Fatalf("Create dish failed, %v.", err)
		} else {
			t.Logf("Dish ID: %v", dish)
		}
	}
}

func TestViewDish(t *testing.T) {
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := "http://localhost:8080/OrcShack/v1/view_dish"

	body := []byte("{ \"id\": 21 }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("View dish failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("View dish failed, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		PrintResponse(res, t)
		t.Fatalf("View dish failed, %v.", err)
	} else {
		err = json.NewDecoder(res.Body).Decode(&dish)
		if err != nil {
			t.Fatalf("View dish failed, %v.", err)
		} else {
			t.Logf("Dish ID: %v", dish)
		}
	}
}

func TestRateDish(t *testing.T) {
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := "http://localhost:8080/OrcShack/v1/rate_dish"

	body := []byte("{ \"id\": 23, \"name\": \"Koni Balls Pizza\" , \"description\": \"Juicy koni balls with blue cheese.\" , \"price\": 8.99, \"category\": \"main course\", \"is_vegetarian\": false,\"is_available\": true,\"rating\": 3.00 }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Update dish failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("Update dish failed, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Update dish failed, %v.", err)
	}
}

func TestDeleteDish(t *testing.T) {
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := "http://localhost:8080/OrcShack/v1/delete_dish"

	body := []byte("{ \"id\": 23  }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Delete dish failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("Delete dish failed, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Delete dish failed, %v.", err)
	} else {
		PrintResponse(res, t)
	}

}
