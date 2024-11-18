package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	sqldb "biz.orcshack/menu/db"
)

// type LoginUser_ResponseUser struct {
// 	JWT string `json:"jwt"`
// }

func List_Dishes(w http.ResponseWriter, r *http.Request) {
	var dish sqldb.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in login request: %v", err)}})
	} else {
	}
}

func View_Dish(w http.ResponseWriter, r *http.Request) {
	var dish sqldb.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in view dish request: %v", err)}})
	} else {
	}
}

func Create_Dish(w http.ResponseWriter, r *http.Request) {
	var dish sqldb.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in create dish request: %v", err)}})
	} else {
		err := sqldb.Global_db.Create_dish(&dish, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Create dish failed: %v", err)}})
		}
	}
}

func Delete_Dish(w http.ResponseWriter, r *http.Request) {
	var dish sqldb.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in delete dish request: %v", err)}})
	} else {
		err := sqldb.Global_db.Delete_dish(dish.Id, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Delete dish failed: %v", err)}})
		}
	}
}

func Search_Dish(w http.ResponseWriter, r *http.Request) {
	var dish sqldb.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in search dish request: %v", err)}})
	} else {
		dishes, err := sqldb.Global_db.Search_dishes_by_name(dish.Name, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Search dish failed: %v", err)}})
		} else {
			// pack the dishes
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*dishes)
		}
	}
}

func Rate_Dish(w http.ResponseWriter, r *http.Request) {
	var dish sqldb.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in rate dish request: %v", err)}})
	} else {
	}
}
