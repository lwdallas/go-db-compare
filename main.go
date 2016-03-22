package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"github.com/gocql/gocql"
	"log"
	"os"
	"strconv"
	"strings"
	// TODO test with csql
	// TODO test with sqlite
	// TODO test with mssql
	// TODO test with psql
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// TODO add windowing of some sort in the source query

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
var firstSession *gocql.Session
var secondSession *gocql.Session

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
	fmt.Println("source Keys:", fmc.SourceKeys)
	fmt.Println("target Keys:", fmc.TargetKeys)
	fmt.Println("ignore field case:", fmc.IgnoreFieldCase)

}

func main() {

	// 0. read config
	fmt.Println("Starting...")

	var err error

	// 1. open source db
	var dbFirst *sql.DB
	var firstCluster *gocql.ClusterConfig
	if sdbc.DatabaseType != "gocql" {
		dbFirst, err = sql.Open(sdbc.DatabaseType, sdbc.ConnectionString)
		if err != nil {
			panic(err.Error())
		}
		defer dbFirst.Close()
	}else {
		firstCluster := gocql.NewCluster(viper.GetString("sourceCluster"))
		firstCluster.Keyspace = viper.GetString("sourceKeyspace")
		firstSession, _ = firstCluster.CreateSession()
		defer firstSession.Close()
	}

	if firstCluster != nil {
		log.Fatal("first cluster setup")
	}

	// 2. open target db

	var dbSecond *sql.DB
	var secondCluster *gocql.ClusterConfig
	if tdbc.DatabaseType != "gocql" {
		dbSecond, err = sql.Open(tdbc.DatabaseType, tdbc.ConnectionString)
		if err != nil {
			panic(err.Error())
		}
		defer dbSecond.Close()
	}else {
		secondCluster := gocql.NewCluster(viper.GetString("targetCluster"))
		secondCluster.Keyspace = viper.GetString("targetKeyspace")
		secondSession, _ = secondCluster.CreateSession()
		if secondCluster == nil {
			log.Fatal("error in first cluster setup")
		}
		defer secondSession.Close()
	}

	if secondCluster != nil {
		log.Fatal("second cluster setup")
	}

	if sdbc.DatabaseType != "gocql" {
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

		rowsCompared := 0

		for rows.Next() {
			// inc count rows
			rowsCompared++

			//  clear result
			for k,_ := range result{
				result[k]=""
			}

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
					log.Println("current row:", rowsCompared, "field ", i, ":", firstResult[i], "!=", secondResult[i])
				}
			}
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		// cassandra version
		iter := firstSession.Query(sourceSql).Iter()

		// get columns
		cols := iter.Columns()
		fmt.Println("Columns()",cols)

		// get rows
		rows,err := iter.RowData()
		if err != nil {
			log.Fatal("RowData() did not work:", err )
		}

		dest, err := iter.SliceMap()
		firstResult := make ([]string, len(dest[0])+1)
		for j := 0; j<len(dest); j++ {
			for k, i := range rows.Columns {
				fmt.Print(k, i)
				if gocql.TypeInt == cols[k].TypeInfo.Type() {
						firstResult[k] = strconv.Itoa(dest[j][i].(int))
				} else if dest[j][i] == nil {
					firstResult[k] = "\\N"
				} else {
					firstResult[k] = dest[j][i].(string)
				}
			}

			err, secondResult := GetColumnsByID(dbSecond, GetIDsFromResultGoCql(firstResult, fmc.SourceKeys, rows))
			if err != nil {
				log.Println("Could not get the Target columns on row ", j)
				log.Fatal(err)
			}

			// compare the result sets

			for i := 0; i < len(firstResult)-1; i++ {
				if firstResult[i] != secondResult[i] {
					// output results
					fmt.Println("current row:", j, "field ", i, ":", firstResult[i], "!=", secondResult[i])
				}
			}
		}

		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println("Done.")
}

func GetColumnsFromARow(rows *sql.Rows, rawResult [][]byte, result []string, dest []interface{}) (error, []string) {
	if tdbc.DatabaseType != "gocql" {
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
	} else {
		// cassandra version
		iter := secondSession.Query(targetSql).Iter()

		// get rows
		rows,err := iter.RowData()
		if err != nil {
			log.Fatal("RowData() did not work:", err )
		}

		dest, err := iter.SliceMap()
		for j := 0; j<len(dest); j++{ // there really should only be one row
			for _, i := range rows.Columns {
				result[j] = dest[j][i].(string)
			}
		}

		if err := iter.Close(); err != nil {
			log.Fatal(err)
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

	result := make([]string, 0)
	if tdbc.DatabaseType != "gocql" {
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
		result = make([]string, len(cols))

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

	} else {
		fmt.Println(aSql)
		// cassandra version
	 	iter := secondSession.Query(aSql).Iter()

		// get columns
		cols := iter.Columns()

		// get rows
		rows,err := iter.RowData()
		if err != nil {
			log.Fatal("RowData() did not work:", err )
		}

		result = make([]string, len(cols))
 		dest, err := iter.SliceMap()
		for j := 0; j<len(dest); j++{ // there really should only be one row
			for k, i := range rows.Columns {
				if gocql.TypeInt == cols[k].TypeInfo.Type() {
					result[k] = strconv.Itoa(dest[j][i].(int))
				} else if dest[j][i] == nil {
					result[k] = "\\N"
				}else{
					result[k] = dest[j][i].(string)
				}
			}
		}

		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}
	return nil, result
}

func GetIDsFromResult(resultSet, keys []string) []string {
	returnIds := make([]string, len(keys))
	for k, v := range keys {
		fmt.Println(k,v)
		returnIds[k] = resultSet[GetCoumnPositionFromSourceName(v)]
	}
	return returnIds
}

func GetIDsFromResultGoCql(resultSet, keys []string, rows gocql.RowData) []string {
	returnIds := make([]string, len(keys))
	for k, v := range keys {
		for i, n := range rows.Columns {
			if v==n {
				returnIds[k] = resultSet[i]
			}
		}
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
