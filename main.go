package main

import (
    "database/sql"
    "fmt"
	"log"
	"os"
	"strconv"
    _ "github.com/go-sql-driver/mysql"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configvar string // config filename

type databaseConfig struct{
	ConnectionString string
	DatabaseType string
}

type fieldConfig struct {
	SourceName string
	TargetName string
	SourceType string
	TargetType string
	Ignore bool
	IgnoreType bool
	IgnoreCase bool
}

type fieldMappingsConfig struct{
	Fields []fieldConfig
	IgnoreFieldCase bool
	SourceKeys []string
	TargetKeys []string
}

var sdbc databaseConfig // source DB config
var tdbc databaseConfig // target DB config
var sourceSql string
var targetSql string

func init() {
	pflag.StringVar(&configvar, "config", "config.json", "configuration file name")
	pflag.Parse()

	fmt.Println("config file", configvar)

	viper.SetConfigFile( configvar )
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
		os.Exit(0)
	}

	err = viper.UnmarshalKey("sourceDatabase", &sdbc)
	if err != nil {
		log.Fatal("unable to decode source DB config into struct, %v", err)
	}
	fmt.Println("source DB:", sdbc)

	err = viper.UnmarshalKey("targetDatabase", &tdbc)
	if err != nil {
		log.Fatal("unable to decode target DB config into struct, %v", err)
	}
	fmt.Println("target DB:", tdbc)

	sourceSql = viper.GetString("sourceSql")
	targetSql = viper.GetString("targetSql")


	var fmc fieldMappingsConfig

	err = viper.UnmarshalKey("fieldMapping", &fmc)
	if err != nil {
		log.Fatal("unable to decode mapping config into struct, %v", err)
	}
	fmt.Println("raw fieldmapping:",fmc.Fields)

	fmt.Println( "configname:", viper.GetString("configName"))
	fmt.Println( "description", viper.GetString("description"))
	fmt.Println( "source table:", viper.GetString("sourceTable"))
	fmt.Println( "target table:", viper.GetString("targetTable"))
	fmt.Println( "source SQL:", viper.GetString("sourceSql"))
	fmt.Println( "target SQL:", viper.GetString("targetSql"))
	fmt.Println( "source Keys:", fmc.SourceKeys[0])
	fmt.Println( "target Keys:", fmc.TargetKeys[0])
	fmt.Println( "ignore field case:", fmc.IgnoreFieldCase)

}

func main() {

	// 0. read config

	// 1. open source db
	dbFirst, err := sql.Open( sdbc.DatabaseType, sdbc.ConnectionString)
    if err != nil {
        panic(err.Error())
    }
    defer dbFirst.Close()

	// 2. open target db

	dbSecond, err := sql.Open( tdbc.DatabaseType, tdbc.ConnectionString)
	if err != nil {
		panic(err.Error())
	}
	defer dbSecond.Close()

    // Execute the source  query
    rows, err := dbFirst.Query( sourceSql )
    if err != nil {
        panic(err.Error())
    }
	defer rows.Close()

	// get the columns for all source records
	cols, err := rows.Columns()
	if err != nil {
		log.Println("Failed to get columns", err)
		return sql.NullString{}
	}
	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		// TODO count rows

		// get the columns
		err = GetColumnsFromARow( row, rawResult, result)
		if err != nil {
			log.Println("Could not get the Source columns")
			log.Fatal(err)
		}

		err = GetColumnsByID( dbSecond, &id)
		if err != nil {
			log.Println("Could not get the Target columns")
			log.Fatal(err)
		}

		// TODO compare the result sets

		// TODO output results
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func GetColumnsFromARow( row *sql.Row, rawResult, result []string ) []string {
	// range thru all rows
	err := row.Scan(dest...)
	if err != nil {
		log.Println("Failed to scan row", err)
		return sql.NullString{}
	}

	for i, raw := range rawResult {
		if raw == nil {
			result[i] = "\\N"
		} else {
			result[i] = string(raw)
		}
	}

	return result

}

func GetColumnsByID( db *sql.DB, id *int) sql.NullString {
	// TODO prepare the query by replacing field tokens

	// Execute the query
	rows, err := db.Query( targetSql + strconv.Itoa(*id))
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println("Failed to get columns", err)
		return sql.NullString{}
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	// range thru all rows
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			log.Println("Failed to scan row", err)
			return sql.NullString{}
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}

		fmt.Printf("%#v\n", result)
	}

//	var (
//		a_id, age int
//		first_name, middle_name, last_name, birthdate, description, more_info, addr, city sql.NullString
//	)
//	for rows.Next() {
//		err := rows.Scan(&a_id, &first_name, &middle_name, &last_name, &age, &birthdate, &description, &more_info, &addr, &city)
//		if err != nil {
//			log.Fatal(err)
//		}
//		return middle_name
//	}
	return sql.NullString{}
}