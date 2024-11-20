package rest_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	sqldb "biz.orcshack/menu/db"
)

//
// Dish API
//

// All roles access, but auth required
func List_Dishes(w http.ResponseWriter, r *http.Request) {
	_, _, err := ValidateJWT(r)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("list dish failed: %v", err)}})
		return
	}

	var dish sqldb.Dish
	err = json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in list dish request: %v", err)}})
	} else {
		dishes, err := sqldb.Global_db.List_dishes_by_category(dish.Category, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("list dish failed: %v", err)}})
		} else {
			// pack the dishes
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*dishes)
		}
	}
}

func View_Dish(w http.ResponseWriter, r *http.Request) {
	_, _, err := ValidateJWT(r)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("view dish failed: %v", err)}})
		return
	}

	var dish sqldb.Dish
	err = json.NewDecoder(r.Body).Decode(&dish)
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

func Search_Dish(w http.ResponseWriter, r *http.Request) {
	_, _, err := ValidateJWT(r)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("search dish failed: %v", err)}})
		return
	}

	var dish sqldb.Dish
	err = json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in search dish request: %v", err)}})
	} else {
		dishes, err := sqldb.Global_db.Search_dishes_by_name(dish.Name, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("search dish failed: %v", err)}})
		} else {
			// pack the dishes
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(*dishes)
		}
	}
}

func Rate_Dish(w http.ResponseWriter, r *http.Request) {
	user_id, _, err := ValidateJWT(r)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Rate dish failed: %v", err)}})
	}

	var udr sqldb.User_Dish_Rating
	err = json.NewDecoder(r.Body).Decode(&udr)
	udr.User_id = user_id
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in rate dish request: %v", err)}})
	} else {
		err := sqldb.Global_db.Create_user_dish_rating(&udr, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("rate dish failed: %v", err)}})
		}
	}
}

// Admin role functions
func Create_Dish(w http.ResponseWriter, r *http.Request) {
	_, role, err := ValidateJWT(r)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("create dish failed: %v", err)}})
		return
	}
	if role == "customer" {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: "create dish failed: user not permitted"}})
		return
	}

	var dish sqldb.Dish
	err = json.NewDecoder(r.Body).Decode(&dish)
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
	_, role, err := ValidateJWT(r)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("delete dish failed: %v", err)}})
		return
	}
	if role == "customer" {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: "delete dish failed: user not permitted"}})
		return
	}

	var dish sqldb.Dish
	err = json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in delete dish request: %v", err)}})
	} else {
		err := sqldb.Global_db.Delete_dish(dish.Id, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("delete dish failed: %v", err)}})
		}
	}
}

// Add a dish image, as many as you like
func Add_Dish_Image(w http.ResponseWriter, r *http.Request) {
	_, role, err := ValidateJWT(r)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("add dish image failed: %v", err)}})
		return
	}
	if role == "customer" {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: "add dish image failed: user not permitted"}})
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(10 << 20) // Limit: 10MB
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the dish ID (as an integer)
	dishIDStr := r.FormValue("dish_id")
	dishID, err := strconv.Atoi(dishIDStr)
	if err != nil {
		http.Error(w, "Invalid dish ID", http.StatusBadRequest)
		return
	}

	// Retrieve the uploaded file
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content into memory
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	// Save the image to the database
	err = sqldb.Global_db.Create_dish_images(uint32(dishID), header.Filename, buf.Bytes(), r.PathValue("restaurant"))
	if err != nil {
		http.Error(w, "Failed to save image to database", http.StatusInternalServerError)
		log.Printf("Error saving image: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Image uploaded successfully"))
}
