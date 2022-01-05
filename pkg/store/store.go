// Package pg will contain the methods to communicate with Postgres
package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PGConfig struct {
	Host               string
	Port               string
	User               string
	Password           string
	DBName             string
	MaxIdleConnections int
	MaxOpenConnections int
	ConnMaxLifetime    time.Duration
}

type PgStorage struct {
	client *sqlx.DB
	psql   sq.StatementBuilderType

	// various SQL queries
	insertUserQuery          string
	isUsernameAvailableQuery string
	isEmailAvailableQuery    string
	isUserValid              string
}

func buildDSN(c *PGConfig) string {
	return fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", c.User, c.Host, c.Port, c.DBName)
}

func Init(c PGConfig) (*PgStorage, error) {
	pg := PgStorage{}

	// Initialize the DB client
	connString := buildDSN(&c)
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		log.Printf("couldn't open db connection: %s\n", err)
		return &pg, err
	}

	log.Println("successful in establishing connection to master db")

	db.SetMaxOpenConns(c.MaxOpenConnections)
	db.SetMaxIdleConns(c.MaxIdleConnections)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		log.Println("couldn't ping master db:", err)
		return &pg, err
	}

	pg.client = db

	pg.psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	pg.insertUserQuery, _, err = pg.psql.
		Insert("users").
		Columns("first_name", "last_name", "username", "email", "created_at", "updated_at").
		Values("$1", "$2", "$3", "$4", "$6", "$7").ToSql()
	if err != nil {
		log.Println("couldn't construct insertUserQuery:", err)
		return &pg, err
	}

	pg.isUsernameAvailableQuery, _, err = pg.psql.
		Select("COUNT(id)").
		From("users").
		Where(sq.Eq{"username": "$1"}).
		ToSql()
	if err != nil {
		log.Println("couldn't construct isUsernameAvailableQuery:", err)
		return &pg, err
	}

	pg.isEmailAvailableQuery, _, err = pg.psql.
		Select("COUNT(id)").
		From("users").
		Where(sq.Eq{"email": "$1"}).
		ToSql()
	if err != nil {
		log.Println("couldn't construct isEmailAvailableQuery:", err)
		return &pg, err
	}

	pg.isUserValid, _, err = pg.psql.
		Select("is_deleted").
		From("users").
		Where(sq.Eq{"user_id": "$1"}).
		ToSql()
	if err != nil {
		log.Println("couldn't construct isUserValid:", err)
		return &pg, err
	}

	return &pg, nil
}

func (pg *PgStorage) Close() {
	err := pg.client.Close()
	if err != nil {
		log.Println("couldn't close db connection:", err)
	}
}

func (pg *PgStorage) InsertUser(firstName, lastName, username, email string) error {
	currentTime := time.Now()
	_, err := pg.client.Exec(pg.insertUserQuery, firstName, lastName, username, email, currentTime, currentTime)
	if err != nil {
		log.Println("couldn't execute InsertUser query:", err)
		return err
	}
	return nil
}

func (pg *PgStorage) IsUsernameAvailable(username string) (bool, error) {
	var count int
	err := pg.client.QueryRow(pg.isUsernameAvailableQuery, username).Scan(&count)
	if err != nil {
		log.Println("couldn't execute IsUsernameAvailable query:", err)
		return false, err
	}
	return count == 0, nil
}

func (pg *PgStorage) IsEmailAvailable(email string) (bool, error) {
	var count int
	err := pg.client.QueryRow(pg.isEmailAvailableQuery, email).Scan(&count)
	if err != nil {
		log.Println("couldn't execute IsEmailAvailable query:", err)
		return false, err
	}
	return count == 0, nil
}

func (pg *PgStorage) IsUserValid(userID int64) (bool, error) {
	var isValid bool
	err := pg.client.QueryRow(pg.isUserValid, userID).Scan(&isValid)
	if err == sql.ErrNoRows {
		return false, err
	}
	if err != nil {
		log.Println("couldn't execute IsUserValid query:", err)
	}

	return isValid, nil
}
