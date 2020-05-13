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

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = resetDB(db, dbname)

	if err != nil {
		panic(err)
	}

	db.Close()

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
