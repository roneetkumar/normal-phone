package db

import (
	"database/sql"
	"fmt"
)

//Phone represent the phone_numbers model
type Phone struct {
	ID     int
	Number string
}

//DB struct
type DB struct {
	db *sql.DB
}

//Open func
func Open(driverName, dataSourse string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourse)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

//Seed func
func (db *DB) Seed() error {
	numbers := []string{"0987654321", "123 231 7890", "(123) 311 7890", "(123) 654 - 7890", "123-981-7890", "(123)234-7890"}

	for _, num := range numbers {
		if _, err := insertPhone(db.db, num); err != nil {
			return err
		}
	}

	return nil
}

// not usefull right now
func getPhone(db *sql.DB, id int) (string, error) {

	var number string

	stm := `SELECT value FROM phone_numbers WHERE id=$1`

	err := db.QueryRow(stm, id).Scan(&number)

	if err != nil {
		return "", err
	}

	return number, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {

	var id int

	stm := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`

	err := db.QueryRow(stm, phone).Scan(&id)

	if err != nil {
		return -1, err
	}

	return int(id), err
}

//Close func
func (db *DB) Close() error {
	return db.db.Close()
}

//FindPhone func
func (db *DB) FindPhone(number string) (*Phone, error) {

	var p Phone

	stm := `SELECT id,value FROM phone_numbers WHERE value=$1`

	err := db.db.QueryRow(stm, number).Scan(&p.ID, &p.Number)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

//UpdatePhone func
func (db *DB) UpdatePhone(p *Phone) error {
	stm := `UPDATE phone_numbers SET value=$2 WHERE id=$1`

	_, err := db.db.Exec(stm, p.ID, p.Number)
	return err
}

//DeletePhone func
func (db *DB) DeletePhone(id int) error {
	stm := `DELETE FROM phone_numbers WHERE id=$1`

	_, err := db.db.Exec(stm, id)
	return err
}

//GetAllPhones func
func (db *DB) GetAllPhones() ([]Phone, error) {

	stm := `SELECT id,value FROM phone_numbers`

	rows, err := db.db.Query(stm)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbNumbers []Phone

	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		dbNumbers = append(dbNumbers, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dbNumbers, err

}

//Migrate func
func Migrate(driverName, dataSourse string) error {

	db, err := sql.Open(driverName, dataSourse)
	if err != nil {
		return err
	}

	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}

	return db.Close()
}

//Reset func
func Reset(driverName, dataSourse, dbName string) error {

	db, err := sql.Open(driverName, dataSourse)
	if err != nil {
		return err
	}

	err = resetDB(db, dbName)
	if err != nil {
		return err
	}

	return db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error {

	stm := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`)

	_, err := db.Exec(stm)
	return err
}

func resetDB(db *sql.DB, name string) error {

	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}

	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {

	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}

	return nil
}
