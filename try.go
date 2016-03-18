package main

import (
    "database/sql"
    "fmt"
	"log"
	"strconv"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dbFirst, err := sql.Open("mysql", "root:password@tcp(:8889)/testgdbc")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer dbFirst.Close()

	dbSecond, err := sql.Open("mysql", "root:password@tcp(:8889)/testgdbc")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer dbSecond.Close()

    // Execute the query
    rows, err := dbFirst.Query("SELECT id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city FROM first")
    if err != nil {
        panic(err.Error())
    }
	defer rows.Close()
	var (
		id, age int
		first_name, middle_name, last_name, birthdate, description, more_info, addr, city sql.NullString
	)
	for rows.Next() {
		err := rows.Scan(&id, &first_name, &middle_name, &last_name, &age, &birthdate, &description, &more_info, &addr, &city)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(id, first_name, middle_name, last_name, birthdate, age, description, more_info, addr, city)
		a_middle_name := GetARowByID( dbSecond, &id)
		if a_middle_name == middle_name{
			fmt.Println(id, "OK")
		} else {
			fmt.Println(id, "target record difference in field middle_name", middle_name, "!=",)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
func GetARowByID( db *sql.DB, id *int) sql.NullString {
	// Execute the query
	rows, err := db.Query("SELECT id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city FROM second WHERE id="+strconv.Itoa(*id))
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var (
		a_id, age int
		first_name, middle_name, last_name, birthdate, description, more_info, addr, city sql.NullString
	)
	for rows.Next() {
		err := rows.Scan(&a_id, &first_name, &middle_name, &last_name, &age, &birthdate, &description, &more_info, &addr, &city)
		if err != nil {
			log.Fatal(err)
		}
		return middle_name
	}
	return sql.NullString{}
}