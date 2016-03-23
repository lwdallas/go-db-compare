# go-db-compare
**Go Database Compare**  
A program that compares data row-by-row between two tables.  Tables can be in different databases, including Cassandra.

The program requires a definition of the databases(source and target)as well as field mappings. The keys and SELECT  statements have to be defined to allow for data to be located. It might be necessary to handle NULLs in the SELECT. Note the LIMIT in the target SELECT.

**Prerequisite Installation**

`go get github.com/gocql/gocql`
`go get github.com/go-sql-driver/mysql`
`go get github.com/spf13/pflag`  
`go get github.com/spf13/viper`  

**Execution**

`go run main.go --config=configuration-file.json`

**Minimum Sample Config File**  
This is an example mysql-to-mysql database compare.  

	{
		"configName": "basic test",
		"description": "This tests the basic row-for row comparison of two mysql tables that are almost identical",
		"sourceDatabase": {
			"connectionString": "root:Password@tcp(127.0.0.1:3306)/testgdbc",
			"connectionStringExplanation": "user:password@tcp(host:portnumber)/database",
			"databaseType": "mysql",
			"databaseTypeExplanation": "mysql, mssql, cql, etc"
		},
		"targetDatabase": {
			"connectionString": "root:Password@tcp(127.0.0.1:3306)/testgdbc",
			"connectionStringExplanation": "user:password@tcp(host:portnumber)/database",
			"databaseType": "mysql",
			"databaseTypeExplanation": "mysql, mssql, cql, etc"
		},
		"sourceTable": "first",
		"targetTable": "second",
		"sourceSql": "SELECT id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city FROM first",
		"targetSql": "SELECT id, first_name, middle_name, last_name, age, birthdate, description, more_info, addr, city FROM second WHERE id={field:id}",
		"fieldMapping": {
			"ignoreFieldCase": true,
			"fields": [{
				"sourceName": "id",
				"targetName": "id"
			}, {
				"sourceName": "first_name",
				"targetName": "first_name"
			}, {
				"sourceName": "middle_name",
				"targetName": "middle_name"
			}, {
				"sourceName": "last_name",
				"targetName": "last_name"
			}, {
				"sourceName": "age",
				"targetName": "age"
			}, {
				"sourceName": "birthdate",
				"targetName": "birthdate"
			}, {
				"sourceName": "description",
				"targetName": "description"
			}, {
				"sourceName": "more_info",
				"targetName": "more_info"
			}, {
				"sourceName": "addr",
				"targetName": "addr"
			}, {
				"sourceName": "city",
				"targetName": "city"
			}],
			"sourceKeys": [
				"id"
			],
			"targetKeys": [
				"id"
			]
		}
	}


**Tested Scenarios** (so far)  
See the test-data directory for sample data and configurations to testing the following scenarios.  

* mysql-mysql  
* gocql-gocql  
* mysql-gocql  
 
**Roadmap**  
Things I need to deal with.  

* I really need to clean up this code base. I just got it  working and checked it in. Reality: This code stinks!
* Add paging to first query to handle potentially large data sets
* Add the ability to compare actual source types with conversion with mapping attribute "sourceType": "varchar"
* Add the ability to compare actual target types with conversion with mapping attribute "targetType": "varchar"
* Add the ability to ignore a field with mapping attribute "ignore": false
* Add the ability to ignore a field's type difference with mapping attribute "ignoreType": true
* Add the ability to ignore a field's case with mapping attribute "ignoreCase": true