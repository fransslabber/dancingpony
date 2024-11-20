package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	sqldb "biz.orcshack/menu/db"
)

func Create_Review(w http.ResponseWriter, r *http.Request) {

	user_id, _, err := ValidateJWT(r)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("create review failed: %v", err)}})
	}

	var review sqldb.Review
	err = json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in create review request: %v", err)}})
	} else {
		err := sqldb.Global_db.Create_review(review.Review, review.Rating, user_id, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("create review failed: %v", err)}})
		}
	}
}
