package sqldb

import (
	"context"
	"time"
)

//
// DB dish review access
//

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

type Array_Reviews []*Review

// Create a review with 0 sentiment score
// Sentiment score updated externally
func (d *SqlDB) Create_review(review string, rating float32, user_id uint32, restaurant string) error {
	_, err := d.db.Exec(context.Background(), "INSERT INTO restaurant_reviews (restaurant_id, user_id, review, rating, sentiment_score)"+
		" VALUES ((select id from restaurants where path_name = $1), $2, $3, $4, 0.00 );", restaurant, user_id, review, rating)
	return err
}

// List reviews by restaurant, return list review array
func (d *SqlDB) List_reviews_by_restaurant(restaurant string) (*Array_Reviews, error) {

	rows, err := d.db.Query(context.Background(), "SELECT * from restaurant_reviews where (select id from restaurants where path_name = $1) = restaurant_id;", restaurant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	reviews := make(Array_Reviews, 0)
	for rows.Next() {
		review := Review{}
		err = rows.Scan(&review.Id, &review.Restaurant_id, &review.User_id, &review.Review, &review.Rating, &review.Sentiment_score, &review.Date_created, &review.Date_updated)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &reviews, nil
}
