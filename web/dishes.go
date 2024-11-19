package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	sqldb "biz.orcshack/menu/db"
)

func List_Dishes(w http.ResponseWriter, r *http.Request) {
	var dish sqldb.Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in list dish request: %v", err)}})
	} else {
		dishes, err := sqldb.Global_db.List_dishes_by_category(dish.Category, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Lis dish failed: %v", err)}})
		} else {
			// pack the dishes
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*dishes)
		}
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
		dish_created, err := sqldb.Global_db.View_dish(dish.Id, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("View dish failed: %v", err)}})
		} else {
			// pack the dishes
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*dish_created)
		}

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
		dish_id, err := sqldb.Global_db.Create_dish(&dish, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Create dish failed: %v", err)}})
		} else {
			dish.Id = dish_id
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(dish)
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
		err := sqldb.Global_db.Update_dish(&dish, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Rate dish failed: %v", err)}})
		}
	}
}
