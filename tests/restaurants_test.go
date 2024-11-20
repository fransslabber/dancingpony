package tests

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type Review struct {
	Id              uint32    `json:"id"`
	Restaurant_id   uint32    `json:"restaurant_id"`
	User_id         uint32    `json:"user_id"`
	Review          string    `json:"review"`
	Rating          float32   `json:"rating"`
	Sentiment_score float32   `json:"sentiment_score"`
	Date_created    time.Time `json:"date_created"`
	Date_updated    time.Time `json:"date_updated"`
}

func TestCreateReview(t *testing.T) {
	teardownSuite, testAuth := SetupLoginTest(t)
	defer teardownSuite(t)

	device_reg_url := "http://localhost:8080/OrcShack/v1/create_review"

	body := []byte("{ \"review\": \"Went on a Monday night thinking we wouldnt have to book - it was full! But they found us a table in 5 mins...the food was great. Massive portions at a VERY good price. Woukd go again and recommend!\" , \"rating\": 4.98 }")

	r, err := http.NewRequest("GET", device_reg_url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Create review failed, %v.", err)
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
		t.Fatalf("Create review failed, %v.", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		PrintResponse(res, t)
		t.Fatalf("Create dish failed, %v.", err)
	}
}
