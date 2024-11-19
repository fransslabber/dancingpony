package sqldb

import (
	"context"
	"fmt"
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

type Array_Dishes []*Dish

func (d *SqlDB) List_dishes_by_category(category, restaurant_id string) (*Array_Dishes, error) {

	rows, err := d.db.Query(context.Background(), "SELECT * from restaurant_dishes where category = $1 AND (select id from restaurants where path_name = $2) = restaurant_id;", category, restaurant_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dishes := make(Array_Dishes, 0)
	for rows.Next() {
		dish := Dish{}
		err = rows.Scan(&dish.Id, &dish.Name, &dish.Description, &dish.Price, &dish.Category, &dish.Is_vegetarian, &dish.Is_available, &dish.Rating, &dish.Restaurant_id, &dish.Date_created, &dish.Date_updated)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, &dish)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &dishes, nil
}

func (d *SqlDB) Search_dishes_by_name(name, restaurant_id string) (*Array_Dishes, error) {

	rows, err := d.db.Query(context.Background(), "SELECT * from restaurant_dishes where name LIKE '%' || $1 || '%' AND (select id from restaurants where path_name = $2) = restaurant_id;", name, restaurant_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dishes := make(Array_Dishes, 0)
	for rows.Next() {
		dish := Dish{}
		err = rows.Scan(&dish.Id, &dish.Name, &dish.Description, &dish.Price, &dish.Category, &dish.Is_vegetarian, &dish.Is_available, &dish.Rating, &dish.Restaurant_id, &dish.Date_created, &dish.Date_updated)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, &dish)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &dishes, nil
}

func (d *SqlDB) Create_dish(dish *Dish, restaurant string) (uint32, error) {
	dish_id := uint32(0)
	err := d.db.QueryRow(context.Background(), "INSERT INTO restaurant_dishes (name, description, price, category, is_vegetarian, is_available, rating, restaurant_id)"+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, (select id from restaurants where path_name = $8)) RETURNING id;",
		dish.Name, dish.Description, dish.Price, dish.Category, dish.Is_vegetarian, dish.Is_available, dish.Rating, restaurant).Scan(&dish_id)
	return dish_id, err
}

func (d *SqlDB) Update_dish(dish *Dish, restaurant string) error {
	_, err := d.db.Exec(context.Background(), "UPDATE restaurant_dishes SET name=$1, description=$2, price=$3, category=$4, is_vegetarian=$5, is_available=$6, rating=$7, restaurant_id=(select id from restaurants where path_name = $8), updated_at = CURRENT_TIMESTAMP where id = $9;",
		dish.Name, dish.Description, dish.Price, dish.Category, dish.Is_vegetarian, dish.Is_available, dish.Rating, restaurant, dish.Id)
	return err
}

func (d *SqlDB) View_dish(id uint32, restaurant string) (*Dish, error) {
	dish := Dish{}
	err := d.db.QueryRow(context.Background(), "SELECT * from restaurant_dishes where id = $1 AND (select id from restaurants where path_name = $2) = restaurant_id;",
		id, restaurant).Scan(&dish.Id, &dish.Name, &dish.Description, &dish.Price, &dish.Category, &dish.Is_vegetarian, &dish.Is_available, &dish.Rating, &dish.Restaurant_id, &dish.Date_created, &dish.Date_updated)
	fmt.Printf("Db view dish %v\n", dish)
	return &dish, err
}

func (d *SqlDB) Delete_dish(id uint32, restaurant string) error {
	_, err := d.db.Exec(context.Background(), "DELETE FROM restaurant_dishes where id = $1 AND (select id from restaurants where path_name = $2) = restaurant_id;", id, restaurant)
	return err
}
