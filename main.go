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

type Phone struct {
	id     int
	number string
}

var numbers []string = []string{"1234567890", "123 456 7890", "(123) 456 7890", "(123) 456 - 7890", "123-456-7890", "(123)456-7890"}

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

	// for _, num := range numbers {
	// 	id, err := insertPhone(db, num)
	// 	must(err)
	// }

	number, err := getPhone(db, 1)

	must(err)

	fmt.Println("Number is ... ", number)

	phones, err := getAllPhones(db)

	must(err)

	for _, p := range phones {
		fmt.Printf("Working on ...%+v\n", p)
		number := normalize(p.number)
		if number != p.number {
			fmt.Println("Updating or removing...", p.number)
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {
				// delete
				must(deletePhone(db, p.id))
			} else {
				//update
				p.number = number
				must(updatePhone(db, p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}

}

func getPhone(db *sql.DB, id int) (string, error) {

	var number string

	stm := `SELECT value FROM phone_numbers WHERE id=$1`

	err := db.QueryRow(stm, id).Scan(&number)

	if err != nil {
		return "", err
	}

	return number, nil
}

func findPhone(db *sql.DB, number string) (*Phone, error) {

	var p Phone

	stm := `SELECT value FROM phone_numbers WHERE value=$1`

	err := db.QueryRow(stm, number).Scan(&p.id, &p.number)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func updatePhone(db *sql.DB, p Phone) error {
	stm := `UPDATE phone_numbers SET value=$2 WHERE id=$1`

	_, err := db.Exec(stm, p.id, p.number)
	return err
}

func deletePhone(db *sql.DB, id int) error {
	stm := `DELETE phone_numbers WHERE id=$1`

	_, err := db.Exec(stm, id)
	return err
}

func getAllPhones(db *sql.DB) ([]Phone, error) {

	stm := `SELECT id,value FROM phone_numbers`

	rows, err := db.Query(stm)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbNumbers []Phone

	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		dbNumbers = append(dbNumbers, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dbNumbers, err

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

func normalize(phone string) string {

	re := regexp.MustCompile("[^0-9]")
	// re := regexp.MustCompile("\\D")

	return re.ReplaceAllString(phone, "")
}

// func normalize(phone string) string {

// 	var buf bytes.Buffer

// 	for _, ch := range phone {

// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}

// 	}
// 	return buf.String()
// }
