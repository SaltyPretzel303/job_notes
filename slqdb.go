package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const CREATE_TABLE_QUERY = `CREATE TABLE Application 
	(	id integer not null primary key, 
		Company text not null,
		Role text, 
		Location text,
		WorkType text, 
		IsHybrid integer,
		IsRemote integer,
		ApplicationDate text,
		ResponseDate text, 
		Response text,
		Comment text
	)`

const CHECK_TABLE_EXISTS_QUERY = `SELECT name FROM sqlite_master 
WHERE type='table' AND name="Application";`

type SqLiteDb struct {
	host string
	conn *sql.DB
}

func (db *SqLiteDb) init() {

	conn, err := sql.Open("sqlite3", db.host)
	db.conn = conn

	if err != nil {
		fmt.Printf("Failed to create SqlLite connection at: %v, err: %v\n", db.host, err)
		return
	}

	rows, err := db.conn.Query(CHECK_TABLE_EXISTS_QUERY)
	if err != nil {
		fmt.Printf("Failed to check if table exists with: %v, err: %v\n",
			CHECK_TABLE_EXISTS_QUERY, err)

		return
	}

	defer rows.Close()

	if !rows.Next() {
		fmt.Println("Will create database using: " + CREATE_TABLE_QUERY)

		_, err = db.conn.Exec(CREATE_TABLE_QUERY)

		if err != nil {
			fmt.Printf("Failed to init SqlLiteDatabase with: %v, err: %v\n", CREATE_TABLE_QUERY, err)
			return
		}

	}
}

func NewLocalDb() *SqLiteDb {
	return NewDb("./sqlite.db")
}

func NewDb(host string) *SqLiteDb {
	db := &SqLiteDb{host: host}
	db.init()
	return db
}

func buildInsertQuery(ap *ApplicationData) string {
	isHybrid, isRemote := 0, 0
	if ap.IsHybrid {
		isHybrid = 1
	}

	if ap.IsRemote {
		isRemote = 1
	}

	return fmt.Sprintf(`INSERT INTO Application 
		(
			Company, Role, Location, WorkType, IsHybrid, IsRemote, 
			ApplicationDate, ResponseDate, Response, Comment
		) 
			VALUES ('%v','%v','%v','%v','%v','%v','%v','%v','%v', '%v')`,
		ap.Company, ap.Role, ap.Location, ap.WorkType, isHybrid, isRemote,
		ap.ApplicationDate, ap.Response, ap.Response, ap.Comment)
}

func (db *SqLiteDb) Save(ap *ApplicationData) {
	fmt.Println("Saving new application in SqlLiteDatabase.")

	insertQuery := buildInsertQuery(ap)
	_, err := db.conn.Exec(insertQuery)

	if err != nil {
		fmt.Println("Failed to exec insert query: ", insertQuery, " err: ", err)
		return
	}

	fmt.Println("Application saved.")
}

func buildSelectSearchQuery(data *SearchData) string {
	conds := []string{}

	dateCond := fmt.Sprintf("ApplicationDate < '%v'", data.OlderThanDate)
	conds = append(conds, dateCond)

	if data.Company != "" {
		query := fmt.Sprintf("Company LIKE \"%%%v%%\"", data.Company)
		conds = append(conds, query)
	}

	if data.Role != "" {
		query := fmt.Sprintf("Role LIKE \"%%%v%%\"", data.Role)
		conds = append(conds, query)
	}

	if data.Location != "" {
		query := fmt.Sprintf("Location LIKE \"%%%v%%\" ", data.Location)
		conds = append(conds, query)
	}

	if data.IsRemote {
		conds = append(conds, "IsRemote=true")
	}

	if data.IsHybrid {
		conds = append(conds, "IsHybrid=true")
	}

	select_part := "SELECT * FROM Application"
	cond_part := strings.Join(conds, " AND ")

	if len(conds) > 0 {
		return fmt.Sprintf("%v WHERE %v", select_part, cond_part)
	} else {
		return select_part
	}
}

func (db *SqLiteDb) Search(searchData *SearchData) []ApplicationData {

	query := buildSelectSearchQuery(searchData)

	// return db.execQuery(query)

	fmt.Println("Executing query: ", query)

	rows, err := db.conn.Query(query)
	if err != nil {
		fmt.Println("Failed search with: ", query)
		return []ApplicationData{}
	}
	defer rows.Close()

	var data = []ApplicationData{}

	for rows.Next() {

		var rowData = ApplicationData{}
		var id int

		err := rows.Scan(&id, &rowData.Company, &rowData.Role, &rowData.Location,
			&rowData.WorkType, &rowData.IsHybrid, &rowData.IsRemote,
			&rowData.ApplicationDate, &rowData.ResponseDate,
			&rowData.Response, &rowData.Comment)

		if err != nil {
			fmt.Println("Failed to read a row: ", err)
		} else {
			data = append(data, rowData)
		}

	}

	return data
}

func (db *SqLiteDb) Close() {
	db.conn.Close()
}
