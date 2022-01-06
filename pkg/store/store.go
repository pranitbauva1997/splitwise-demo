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
	allUsersQuery            string
	addBillQuery             string
	addTransactionQuery      string
	userTransactionsQuery    string
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
		Values("$1", "$2", "$3", "$4", "$5", "$6").ToSql()
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

	pg.allUsersQuery, _, err = pg.psql.
		Select("id", "first_name", "last_name", "username", "email").
		From("users").
		Where(sq.Eq{"is_deleted": "$1"}).ToSql()
	if err != nil {
		log.Println("couldn't construct isUserValid:", err)
		return &pg, err
	}
	pg.addBillQuery, _, err = pg.psql.
		Insert("bills").
		Columns("user_id", "amount", "created_at", "updated_at", "is_deleted").
		Values("$1", "$2", "$3", "$4", "$5").Suffix("RETURNING id").ToSql()
	if err != nil {
		log.Println("couldn't construct addBillQuery:", err)
		return &pg, err
	}
	pg.addTransactionQuery, _, err = pg.psql.
		Insert("transactions").
		Columns("bill_id", "owed_to", "owes", "amount", "created_at", "updated_at", "is_deleted").
		Values("$1", "$2", "$3", "$4", "$5", "$6", "$7").ToSql()
	if err != nil {
		log.Println("couldn't construct addTransactionQuery:", err)
		return &pg, err
	}
	pg.userTransactionsQuery, _, err = pg.psql.
		Select("owed_to", "owes", "amount").
		From("transactions").
		Where("(owed_to=$1 OR owes=$1) AND (owed_to != owes) AND is_deleted=false;").ToSql()
	if err != nil {
		log.Println("couldn't construct userTransactionsQuery:", err)
		return &pg, err
	}
	log.Println(pg.userTransactionsQuery)

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
		return fmt.Errorf("couldn't execute InsertUser query: %w", err)
	}
	return nil
}

func (pg *PgStorage) IsUsernameAvailable(username string) (bool, error) {
	var count int
	err := pg.client.QueryRow(pg.isUsernameAvailableQuery, username).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("couldn't execute IsUsernameAvailable query: %w", err)
	}
	return count == 0, nil
}

func (pg *PgStorage) IsEmailAvailable(email string) (bool, error) {
	var count int
	err := pg.client.QueryRow(pg.isEmailAvailableQuery, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("couldn't execute IsEmailAvailable query: %w", err)
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
		return false, fmt.Errorf("couldn't execute IsUserValid query: %w", err)
	}

	return isValid, nil
}

func (pg *PgStorage) GetAllUsers() ([]UserResponse, error) {
	allUsers := make([]UserResponse, 0)
	rows, err := pg.client.Query(pg.allUsersQuery, false)
	if err != nil {
		return nil, fmt.Errorf("couldn't execute GetAllUsers query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user UserResponse
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.UserName, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse data returned by GetAllUsers query: %w", err)
		}
		allUsers = append(allUsers, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("coudn't read the data returned by GetAllUsers query: %w", err)
	}

	return allUsers, nil
}

func (pg *PgStorage) AddBill(bill AddBillRequest) error {
	currentTime := time.Now()
	var billId int64

	tx, err := pg.client.Beginx()
	if err != nil {
		return fmt.Errorf("couldn't start the DB transaction for AddBill query: %w", err)
	}
	err = tx.QueryRow(pg.addBillQuery, bill.CreatedBy, bill.Amount, currentTime, currentTime, false).Scan(&billId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("couldn't execute AddBill query: %w", err)
	}

	for _, t := range bill.Transactions {
		_, err = pg.client.Exec(pg.addTransactionQuery, billId, t.OwedTo, t.Owes, t.Amount, currentTime, currentTime, false)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("couldn't execute AddTransaction query: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("couldn't commit the transaction for AddBill query: %w", err)
	}

	return nil
}

func (pg *PgStorage) UserTransactions(userId int64) ([]UserTransaction, error) {
	userTransactions := make([]UserTransaction, 0)
	rows, err := pg.client.Query(pg.userTransactionsQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("couldn't execute UserTransactions query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var userTransaction UserTransaction
		err = rows.Scan(&userTransaction.OwedTo, &userTransaction.Owes, &userTransaction.Amount)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse data returned by UserTransactions query: %w", err)
		}
		userTransactions = append(userTransactions, userTransaction)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("coudn't read the data returned by UserTransactions query: %w", err)
	}

	return userTransactions, nil
}