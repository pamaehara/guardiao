package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq"
)

const (
	host     = "rds-gwapi-cvc-dev.reservafacil.tur.br"
	port     = 5432
	user     = "gwapi"
	password = "aFTwL7aZk7T3g"
	dbname   = "gwapi"
)

func psqlInfoDev() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

// GetConn abre conexão com o DB
func GetConn() *gorm.DB {
	// db, err := sql.Open("postgres", psqlInfoDev())
	db, err := gorm.Open("postgres", psqlInfoDev())
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	fmt.Println("Successfully connected!")
	return db
}
