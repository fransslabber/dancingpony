package sqldb

import (
	"context"
	"time"
)

// User DB access
type User struct {
	Id            uint32
	Name          string `json:"name"`
	Email         string `json:"email"`
	Role          string
	Restaurant_id uint32
	Password_hash string `json:"password"`
	Salt          string
	Date_created  time.Time
	Date_updated  time.Time
}

type Array_Users []*User

// Register a new user via api - role can only customer
// For other roles need admin restricted access
// Hashed password in DB, has separate salt column - used depending on hash algo
func (d *SqlDB) Create_user(name string, email string, password string, restaurant string) error {
	_, err := d.db.Exec(context.Background(), "INSERT INTO users (name, email, role, restaurant_id, hashed_password, salt)"+
		" VALUES ($1, $2, 'customer',(select id from restaurants where path_name = $3), crypt($4, gen_salt('bf')), '');", name, email, restaurant, password)
	return err
}

// Authenticate user, return true/false and user details
func (d *SqlDB) Login_user(email string, password string, restaurant string) (bool, *User, error) {
	var is_authenticated bool
	user := User{}
	err := d.db.QueryRow(context.Background(), "SELECT (hashed_password = crypt($3, hashed_password)) AS is_authenticated,* FROM users WHERE email = $1 AND (select id from restaurants where path_name = $2) = restaurant_id;",
		email, restaurant, password).Scan(&is_authenticated, &user.Id, &user.Name, &user.Email, &user.Role, &user.Restaurant_id, &user.Password_hash, &user.Salt, &user.Date_created, &user.Date_updated)
	return is_authenticated, &user, err
}
