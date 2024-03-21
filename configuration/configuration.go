package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func LoadDatabase(filename string) (*Database, error) {
	//reading the file
	filepath, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error in reading the file", err)
		return nil, err
	}

	//unmarshal the contents
	var Database_connections Database
	err = json.Unmarshal(filepath, &Database_connections)
	fmt.Println(&Database_connections)
	return &Database_connections, err
}
