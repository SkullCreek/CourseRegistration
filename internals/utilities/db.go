package utilities

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// parameters required to connect with database
type connection struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

// This method is used to setup connect with database
func InitDB() *sql.DB {
	conn := connection{
		Host:     "mouse.db.elephantsql.com",
		Port:     "5432",
		User:     "jjnyldpb",
		Password: "8SlAuwZXe2RnHbtjkshuFheLrkPPmETs",
		Dbname:   "jjnyldpb",
	}
	db, err := sql.Open("postgres", connToString(conn))
	HandleError(err)
	ping_err := db.Ping()
	HandleError(ping_err)
	return db

}

// Converts the connection interface to a string
func connToString(conn connection) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", conn.Host, conn.Port, conn.User, conn.Password, conn.Dbname)
}
