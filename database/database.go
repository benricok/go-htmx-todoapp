package database

import (
	"database/sql"
	"fmt"
	"go-htmx/user"
	"log"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

//type Database struct {
var	Db *sql.DB

func NewDatabase(url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS users (
				username TEXT NOT NULL,
				password TEXT NOT NULL,
				name TEXT,
				surname TEXT,
				email TEXT,
				datecreated TIMESTAMPTZ,
				PRIMARY KEY (username)
			);
			CREATE TABLE IF NOT EXISTS items (
				id TEXT NOT NULL,
				owner TEXT NOT NULL,
				name TEXT NOT NULL,
				done INTEGER NOT NULL,
				datecreated TIMESTAMPTZ NOT NULL,
				datecompleted TIMESTAMPTZ,
				PRIMARY KEY (id),
				FOREIGN KEY (owner) REFERENCES users(username) ON DELETE CASCADE
			);
	`)

	Db = db
	return nil
}

func GetUser(username string) (*user.User,  error) {
	res, err := Db.Query(`SELECT * FROM users WHERE username = $1`, username)
	if err != nil {
		return nil, fmt.Errorf("unable to scan db row: %+v", err)
	}

	res.Next()
	var (
		password string
		name string
		surname string
		email string
		datecreated time.Time
	)

    err = res.Scan(&username, &password, &name, &surname, &email, &datecreated)
    if err != nil {
        return nil, fmt.Errorf("unable to scan db row: %+v", err)
    }

    res.Close()

    return &user.User{
        Username: username,
        Password: password,
        Name: name,
        Surname: surname,
		Email: email,
        Created: datecreated,
    }, nil
}

func AddUser(u *user.User) (error) {
	res, err := Db.Exec(`INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6)`, u.Username,
										u.Password, u.Name, u.Surname, u.Email, u.Created)
	if err != nil {
		return err
	}
									
	rowsAffected, err:= res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Added user %s, %d rows affected", u.Username, rowsAffected)
	return nil
}

func DeleteUser(username string) (error) {
	res, err := Db.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		return err
	}

	rowsAffected, err:= res.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Deleted user %s, %d rows affected", username, rowsAffected)
	return nil
}