package tests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
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

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
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

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
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

	body := []byte("{ \"id\": 24 }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("View dish failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
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

	body := []byte("{ \"dish_id\": 7, \"rating\": 2.15 }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Update dish failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("Update dish failed, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		PrintResponse(res, t)
		t.Fatalf("Update dish failed, %v.", err)
	}
}

func TestDeleteDish(t *testing.T) {
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := "http://localhost:8080/OrcShack/v1/delete_dish"

	body := []byte("{ \"id\": 25  }")
	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Delete dish failed, %v.", err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("Delete dish failed, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		PrintResponse(res, t)
		t.Fatalf("Delete dish failed, %v.", err)
	} else {
		PrintResponse(res, t)
	}

}

func TestAddDishImage(t *testing.T) {
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := "http://localhost:8080/OrcShack/v1/add_dish_image"

	// Open the image file
	file, err := os.Open("../assets/mouse.jpg")
	if err != nil {
		t.Fatalf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Create a buffer to hold the multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the user_id field
	err = writer.WriteField("dish_id", "7")
	if err != nil {
		t.Fatalf("failed to write dish field: %v", err)
	}

	// Add the image file field
	part, err := writer.CreateFormFile("image", file.Name())
	if err != nil {
		t.Fatalf("failed to create form file field: %v", err)
	}

	// Copy the file content into the form file field
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("failed to copy file content: %v", err)
	}

	// Close the multipart writer to finalize the form
	err = writer.Close()
	if err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", device_reg_url, body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", testAuth.Session_jwt))

	// Send the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Bypass SSL verification (not recommended for production)
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		PrintResponse(resp, t)
		t.Fatalf("add dish image failed, %v.", err)
	}
}
