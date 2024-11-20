package sqldb

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB Connection Config
const (
	host     = "localhost"
	port     = 5432
	user     = "dancingponysvc"
	password = "password"
	dbname   = "dancingpony"
)

type SqlDB struct {
	db *pgxpool.Pool
}

var Global_db SqlDB

func CreateDB() error {
	err := Global_db.Open("")
	if err != nil {
		log.Printf("Could not open pg connections pool: %+v", err)
		return err
	}
	return nil
}

func CloseDB() {
	Global_db.Close()
}

func (s *SqlDB) Open(file string) error {
	var err error
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	s.db, err = pgxpool.New(context.Background(), psqlconn)
	return err
}

func (s *SqlDB) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
