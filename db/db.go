package sqldb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (s *SqlDB) Open(file string) error {
	var err error
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	s.db, err = pgxpool.New(context.Background(), psqlconn)
	return err
}

func (s *SqlDB) Build_version(major uint32, minor uint32, revision uint32) uint32 {
	// Current system version
	nversion := revision
	nversion |= minor << 8
	nversion |= major << 16
	return nversion

}

func (s *SqlDB) Unbuild_version(version uint32) string {
	return fmt.Sprintf("%d.%d.%d", (version>>16)&0xFF, (version>>8)&0xFF, version&0xFF)
}

func (s *SqlDB) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
