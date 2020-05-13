package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "abc123..."
	dbname   = "phone"
)

func main() {

	fmt.Println("Phone Number Normalizer")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	// db, err := sql.Open("postgres", psqlInfo)
	// must(err)

	// err = resetDB(db, dbname)
	// must(err)
	// db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(createPhoneNumbersTable(db))

	id, err := insertPhone(db, "0191561212")

	must(err)

	fmt.Println(id)

}

func insertPhone(db *sql.DB, phone string) (int, error) {

	stm := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`

	var id int

	err := db.QueryRow(stm, phone).Scan(&id)

	if err != nil {
		return -1, err
	}

	return int(id), err
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

func must(err error) {
	if err != nil {
		panic(err)
	}
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

func normalizer(phone string) string {

	re := regexp.MustCompile("[^0-9]")
	// re := regexp.MustCompile("\\D")

	return re.ReplaceAllString(phone, "")
}

// func normalizer(phone string) string {

// 	var buf bytes.Buffer

// 	for _, ch := range phone {

// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}

// 	}
// 	return buf.String()
// }
