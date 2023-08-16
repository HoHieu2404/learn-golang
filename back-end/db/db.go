package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type DBI struct {
	db *sql.DB
}

type DBIInterface interface {
	InitDB() error
	GetDatabase() *sql.DB
}

func tableExists(db *sql.DB, dbName string, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?")
	err := db.QueryRow(query, dbName, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (dbi *DBI) InitDB() error {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env")
	}
	dbHost := os.Getenv("MYSQL_HOST")
	dbPORT := os.Getenv("MYSQL_PORT")
	dbNAME := os.Getenv("MYSQL_DATABASE")
	dbUSER := os.Getenv("MYSQL_USER")
	dbPASSWORD := os.Getenv("MYSQL_PASSWORD")

	url := dbUSER + ":" + dbPASSWORD + "@tcp(" + dbHost + ":" + dbPORT + ")/" + dbNAME

	db, _ := sql.Open("mysql", url)
	err = db.Ping()
	if err != nil {
		panic("Connect to database: " + err.Error())
	}
	dbi.db = db

	exists, err := tableExists(db, dbNAME, "Rates")
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err := db.Exec(`CREATE TABLE Rates (
			id INT PRIMARY KEY AUTO_INCREMENT,
			date DATE,
			currency VARCHAR(255),
			rate FLOAT
		)`)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Rates table has been created")
	}
	return nil
}

func (dbi *DBI) GetDatabase() *sql.DB {
	return dbi.db
}

func NewDBI() DBIInterface {
	return &DBI{}
}
