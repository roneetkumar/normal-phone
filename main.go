package main

import (
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
	pkgdb "github.com/roneetkumar/normal-phone/db"
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

	//Reseting the db
	must(pkgdb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)

	must(pkgdb.Migrate("postgres", psqlInfo))

	db, err := pkgdb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	err = db.Seed()
	must(err)

	phones, err := db.GetAllPhones()
	must(err)

	for _, p := range phones {
		fmt.Printf("Working on ...%+v\n", p)
		number := normalize(p.Number)

		if number != p.Number {
			fmt.Println("Updating or removing...", p.Number)
			existing, err := db.FindPhone(number)
			must(err)

			if existing != nil {
				// delete
				must(db.DeletePhone(p.ID))
			} else {
				//update
				p.Number = number
				must(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}

}

func must(err error) {
	if err != nil {
		panic(err)
	}
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
