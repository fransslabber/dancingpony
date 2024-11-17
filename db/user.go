package sqldb

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	Id            uint32
	Tenant_id     uint32
	Name          string
	Email         string
	Password_hash []byte
	Salt          string
	Last_login    sql.NullTime
	date_created  time.Time
	date_updated  time.Time
}

type Array_Users []*User

func (d *SqlDB) Load_users_by_user_name(user_name string) (error, *Array_Users) {
	return d.load_users_sql("select * from [dbo].[user] where username = '" + user_name + "';")
}

// func (d *SqlDB) Load_users_by_user_id(user_id uint32) (error, *Array_Users) {
// 	return d.load_users_sql("select * from [dbo].[user] where user_id = '" + strconv.FormatUint(uint64(user_id), 10) + "';")
// }

// func (d *SqlDB) Load_users() (error, *Array_Users) {
// 	return d.load_users_sql("select * from [dbo].[user];")
// }

func (d *SqlDB) load_users_sql(sql string) (error, *Array_Users) {
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

func (d *SqlDB) Create_user_sp(name string, user_name string, passwd string, card_id string, auth_type string) (uint32, error) {
	user_id := uint32(0)
	_, err := d.db.Exec("spAdduser", sql.Named("pName", name), sql.Named("pUserName", user_name), sql.Named("pPassword", passwd), sql.Named("pCardId", card_id), sql.Named("pAuthType", auth_type), sql.Named("UserID", sql.Out{Dest: &user_id}))
	if err != nil {
		return 0, err
	}
	if user_id == 0 {
		return 0, fmt.Errorf("User already exists.")
	}
	return user_id, nil
}

func (d *SqlDB) Login_user_sp(user string, password string, device_uuid string) (string, uint32, error) {
	account := "User successfully logged in"
	userid := 0
	_, err := d.db.Exec("spUserLogin", sql.Named("pUserName", user), sql.Named("pPassword", password), sql.Named("pDeviceUUID", device_uuid), sql.Named("responseMessage", sql.Out{Dest: &account}), sql.Named("UserID", sql.Out{Dest: &userid}))
	return account, uint32(userid), err
}

func (d *SqlDB) Login_card_sp(card string) (string, uint32, error) {
	account := "User successfully logged in"
	userid := 0
	_, err := d.db.Exec("spUserLogin", sql.Named("pCard", card), sql.Named("responseMessage", sql.Out{Dest: &account}), sql.Named("UserID", sql.Out{Dest: &userid}))
	return account, uint32(userid), err
}

// func (d *SqlDB) Delete_user(user_name string) (bool, uint32) {

// 	res, err := d.db.Exec("delete from [dbo].[user] where user_name = '" + user_name + "'")
// 	if err != nil {
// 		d.errmsg = err.Error()
// 		return false, 0
// 	}

// 	id, err := res.RowsAffected()
// 	if err != nil {
// 		d.errmsg = err.Error()
// 		return false, 0
// 	}
// 	return true, uint32(id)
// }

// func (d *SqlDB) Create_user_history(user_id uint32, client_id uint32, login_status bool) (error, uint32) {

// 	var rows *sql.Rows
// 	var err error
// 	if login_status {
// 		rows, err = d.db.Query("insert into user_history (user_id, client_id, action, date_created ) OUTPUT INSERTED.[user_history_id] values( " +
// 			strconv.FormatUint(uint64(user_id), 10) + ", " + strconv.FormatUint(uint64(client_id), 10) + " , 'LOGIN',  SYSDATETIMEOFFSET());")
// 	} else {
// 		rows, err = d.db.Query("insert into user_history (user_id, client_id, action, date_created ) OUTPUT INSERTED.[user_history_id] values( " +
// 			strconv.FormatUint(uint64(user_id), 10) + ", " + strconv.FormatUint(uint64(client_id), 10) + " , 'LOGOUT',  SYSDATETIMEOFFSET());")
// 	}

// 	if err != nil {
// 		return err, 0
// 	}

// 	if rows.Next() {
// 		user_history_id := uint32(0)
// 		err = rows.Scan(&user_history_id)
// 		if err != nil {
// 			return err, 0
// 		}
// 		rows.Close()
// 		return nil, user_history_id
// 	}
// 	if err = rows.Err(); err != nil {
// 		return err, 0
// 	}
// 	return fmt.Errorf("Create user history failed with unknown error."), 0
// }

// func (d *SqlDB) Create_audit_log(user_id uint32, client_id uint32, action string, detail string) (error, uint32) {

// 	rows, err := d.db.Query("insert into audit_log (user_id, client_id, action, detail, date_created ) OUTPUT INSERTED.[audit_log_id] values( " +
// 		strconv.FormatUint(uint64(user_id), 10) + ", " + strconv.FormatUint(uint64(client_id), 10) + " , '" + action + "', '" + detail + "' , SYSDATETIMEOFFSET());")

// 	if err != nil {
// 		return err, 0
// 	}

// 	if rows.Next() {
// 		audit_log_id := uint32(0)
// 		err = rows.Scan(&audit_log_id)
// 		rows.Close()
// 		return nil, audit_log_id
// 	}
// 	if err = rows.Err(); err != nil {
// 		return err, 0
// 	}
// 	return fmt.Errorf("Create audit log failed with unknown error."), 0
// }

type Group struct {
	Id           uint32
	Name         string
	Parent_id    uint32
	Filter       string
	Date_created time.Time
}
type Array_Groups []*Group

func (d *SqlDB) Create_group(grp *Group) (uint32, error) {

	parent_id := "NULL"
	if grp.Parent_id != 0 {
		parent_id = fmt.Sprintf("%d", grp.Parent_id)
	}
	insertStr := fmt.Sprintf("insert into [dbo].[group] (name,parent_group_id,filter,date_created) OUTPUT INSERTED.[id] values( '%s',%s, '%s', SYSDATETIMEOFFSET());",
		grp.Name, parent_id, grp.Filter)

	rows, err := d.db.Query(insertStr)
	if err != nil {
		return 0, err
	}

	if rows.Next() {
		defer rows.Close()
		group_id := uint32(0)
		err = rows.Scan(&group_id)
		return group_id, err
	}
	return 0, err
}

func (d *SqlDB) Delete_groups_by_name(name string) (uint32, error) {

	res, err := d.db.Exec("delete from [dbo].[group] where name = '" + name + "'")
	if err != nil {
		return 0, err
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// func (d *SqlDB) Load_groups_by_name(user_name string) (error, *Array_Users) {
// 	return d.load_users_sql("select * from [dbo].[user] where username = '" + user_name + "';")
// }

// func (d *SqlDB) Load_users_by_user_id(user_id uint32) (error, *Array_Users) {
// 	return d.load_users_sql("select * from [dbo].[user] where user_id = '" + strconv.FormatUint(uint64(user_id), 10) + "';")
// }

// func (d *SqlDB) Load_users() (error, *Array_Users) {
// 	return d.load_users_sql("select * from [dbo].[user];")
// }

// func (d *SqlDB) load_users_sql(sql string) (error, *Array_Users) {
// 	rows, err := d.db.Query(sql)
// 	if err != nil {
// 		return err, nil
// 	}
// 	defer rows.Close()
// 	users := make(Array_Users, 0)
// 	for rows.Next() {
// 		usr := User{}
// 		err = rows.Scan(&usr.User_id, &usr.Name, &usr.Username, &usr.Password_hash, &usr.Salt, &usr.Card_id, &usr.Auth_type, &usr.Last_login, &usr.date_created, &usr.date_updated)
// 		if err != nil {
// 			return err, nil
// 		}
// 		users = append(users, &usr)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return err, nil
// 	}
// 	return nil, &users
