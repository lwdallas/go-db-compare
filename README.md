# go-db-compare
Go Database Compare


Installation

`go get github.com/go-sql-driver/mysql`  
`go get github.com/spf13/pflag`  
`go get github.com/spf13/viper`  

Execution

`go run main.go --config=configuration-file.json`

Sample Config File

	{
		"configName": "basic test",
		"description": "This tests the basic row-for row comparison of two mysql tables that are almost identical",
		"sourceDatabase": {
			"connectionString": "root:Password@tcp(127.0.0.1:8889)/testgdbc",
			"connectionStringExplanation": "user:password@tcp(host:portnumber)/database",
			"databaseType": "mysql",
			"databaseTypeExplanation": "mysql, mssql, cql, etc"
		},
		"targetDatabase": {
			"connectionString": "root:Password@tcp(127.0.0.1:8889)/testgdbc",
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
				"targetName": "id",
				"sourceType": "int",
				"targetType": "int",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "first_name",
				"targetName": "first_name",
				"sourceType": "varchar",
				"targetType": "varchar",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "middle_name",
				"targetName": "middle_name",
				"sourceType": "varchar",
				"targetType": "varchar",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "last_name",
				"targetName": "last_name",
				"sourceType": "varchar",
				"targetType": "varchar",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "age",
				"targetName": "age",
				"sourceType": "int",
				"targetType": "int",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "birthdate",
				"targetName": "birthdate",
				"sourceType": "varchar",
				"targetType": "varchar",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "description",
				"targetName": "description",
				"sourceType": "varchar",
				"targetType": "varchar",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "more_info",
				"targetName": "more_info",
				"sourceType": "varchar",
				"targetType": "varchar",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "addr",
				"targetName": "addr",
				"sourceType": "varchar",
				"targetType": "varchar",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}, {
				"sourceName": "city",
				"targetName": "city",
				"sourceType": "varchar",
				"targetType": "varchar",
				"ignore": false,
				"ignoreType": true,
				"ignoreCase": true
			}],
			"sourceKeys": [
				"id"
			],
			"targetKeys": [
				"id"
			]
		}
	}


