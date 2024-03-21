package controller

import (
	"careercraftsman_chatbot/configuration"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connecting_database() (*sql.DB, error) {
	filepath, err := configuration.LoadDatabase("config.json")
	if err != nil {
		log.Println("Error in loading the file")
	}

	dbconfig := configuration.Database{
		Host:     filepath.Host,
		Port:     filepath.Port,
		User:     filepath.User,
		Password: filepath.Password,
		Dbname:   filepath.Dbname,
	}

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable", dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Dbname)
	fmt.Println(psql)

	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.Println("Error in Connecting to database", err)
		return nil, err
	}

	//defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Error in pinging databse", err)
		return nil, err
	}
	log.Println("Sucessfully connected to database")
	return db, err
}
