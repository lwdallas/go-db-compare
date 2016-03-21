package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	// TODO test with csql
	// TODO test with sqlite
	// TODO test with mssql
	// TODO test with psql
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configvar string // config filename

type databaseConfig struct {
	ConnectionString string
	DatabaseType     string
}

type fieldConfig struct {
	SourceName string
	TargetName string
	SourceType string
	TargetType string
	Ignore     bool
	IgnoreType bool
	IgnoreCase bool
}

type fieldMappingsConfig struct {
	Fields          []fieldConfig
	IgnoreFieldCase bool
	SourceKeys      []string
	TargetKeys      []string
}

var sdbc databaseConfig // source DB config
var tdbc databaseConfig // target DB config
var sourceSql string
var targetSql string
var fmc fieldMappingsConfig

func init() {
	pflag.StringVar(&configvar, "config", "config.json", "configuration file name")
	pflag.Parse()

	fmt.Println("config file", configvar)

	viper.SetConfigFile(configvar)
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

	err = viper.UnmarshalKey("fieldMapping", &fmc)
	if err != nil {
		log.Fatal("unable to decode mapping config into struct, %v", err)
	}
	fmt.Println("raw fieldmapping:", fmc.Fields)

	fmt.Println("configname:", viper.GetString("configName"))
	fmt.Println("description", viper.GetString("description"))
	fmt.Println("source table:", viper.GetString("sourceTable"))
	fmt.Println("target table:", viper.GetString("targetTable"))
	fmt.Println("source SQL:", viper.GetString("sourceSql"))
	fmt.Println("target SQL:", viper.GetString("targetSql"))
	fmt.Println("source Keys:", fmc.SourceKeys[0])
	fmt.Println("target Keys:", fmc.TargetKeys[0])
	fmt.Println("ignore field case:", fmc.IgnoreFieldCase)

}

func main() {

	// 0. read config

	// 1. open source db
	dbFirst, err := sql.Open(sdbc.DatabaseType, sdbc.ConnectionString)
	if err != nil {
		panic(err.Error())
	}
	defer dbFirst.Close()

	// 2. open target db

	dbSecond, err := sql.Open(tdbc.DatabaseType, tdbc.ConnectionString)
	if err != nil {
		panic(err.Error())
	}
	defer dbSecond.Close()

	// 3. Execute the source  query
	rows, err := dbFirst.Query(sourceSql)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// 4. get the columns for all source records
	cols, err := rows.Columns()
	if err != nil {
		log.Println("Failed to get columns", err)
		return
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
		//  clear result

		// get the columns
		err, firstResult := GetColumnsFromARow(rows, rawResult, result, dest)
		if err != nil {
			log.Println("Could not get the Source columns")
			log.Fatal(err)
		}

		err, secondResult := GetColumnsByID(dbSecond, GetIDsFromResult(firstResult, fmc.SourceKeys))
		if err != nil {
			log.Println("Could not get the Target columns")
			log.Fatal(err)
		}

		// compare the result sets

		for i := 0; i < len(firstResult); i++ {
			if firstResult[i] != secondResult[i] {
				// output results
				log.Println("<<current row>>:", "field ", i, ":", firstResult[i], "!=", secondResult[i])
			}
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func GetColumnsFromARow(rows *sql.Rows, rawResult [][]byte, result []string, dest []interface{}) (error, []string) {
	// range thru all rows
	err := rows.Scan(dest...)
	if err != nil {
		log.Println("Failed to scan row", err)
		return err, result
	}

	for i, raw := range rawResult {
		if raw == nil {
			result[i] = "\\N"
		} else {
			result[i] = string(raw)
		}
	}

	return nil, result

}

func GetColumnsByID(db *sql.DB, ids []string) (error, []string) {
	aSql := targetSql
	// modify the query by replacing field tokens
	for k, v := range ids {
		aSql = ReplaceToken(aSql, fmc.SourceKeys[k], v)
	}

	// Execute the query
	rows, err := db.Query(aSql)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println("Failed to get columns", err)
		return err, nil
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
			return err, result
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
	return nil, result
}

func GetIDsFromResult(resultSet, keys []string) []string {
	returnIds := make([]string, len(keys))
	for k, v := range keys {
		returnIds[k] = resultSet[GetCoumnPositionFromSourceName(v)]
	}
	return returnIds
}

func GetPositionsFromResult(resultSet, keys []string) []int {
	returnIds := make([]int, len(keys))
	for k, v := range keys {
		returnIds[k] = GetCoumnPositionFromSourceName(v)
	}
	return returnIds
}

func GetCoumnPositionFromSourceName(sourceName string) int {
	for i := 0; i < len(fmc.Fields); i++ {
		if fmc.Fields[i].SourceName == sourceName {
			return i
		}
	}
	return -1
}

func ReplaceToken(aString string, fieldName string, value string) string {
	token := "{field:" + fieldName
	loc := strings.Index(aString, token)
	if loc == -1 {
		return aString
	}
	s := aString[0:loc]
	s2 := aString[loc+len(token):]
	loc = strings.Index(s2, "}")
	s2 = s2[loc+1:]
	return s + value + s2
}
