package sqldb

import (
	"time"
)

type Tenant struct {
	Id           uint32
	Name         string
	date_created time.Time
	date_updated time.Time
}

type Array_Tenants []*Tenant

func (d *SqlDB) Load_tenants() (error, *Array_Users) {
	return d.load_users_sql("select * from [dbo].[user] where username = '" + user_name + "';")
}

// func (d *SqlDB) Load_users_by_user_id(user_id uint32) (error, *Array_Users) {
// 	return d.load_users_sql("select * from [dbo].[user] where user_id = '" + strconv.FormatUint(uint64(user_id), 10) + "';")
// }

// func (d *SqlDB) Load_users() (error, *Array_Users) {
// 	return d.load_users_sql("select * from [dbo].[user];")
// }

func (d *SqlDB) load_tenants_sql(sql string) (error, *Array_Users) {
	rows, err := d.db.Query(sql)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	users := make(Array_Users, 0)
	for rows.Next() {
		usr := User{}
		err = rows.Scan(&usr.User_id, &usr.Name, &usr.Username, &usr.Password_hash, &usr.Salt, &usr.Card_id, &usr.Auth_type, &usr.Last_login, &usr.date_created, &usr.date_updated)
		if err != nil {
			return err, nil
		}
		users = append(users, &usr)
	}
	if err = rows.Err(); err != nil {
		return err, nil
	}
	return nil, &users
}
