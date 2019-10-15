package controllers

import (
	"fmt"
	"github.com/davidchandra95/fullstack/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

type Server struct {
	DB *gorm.DB
	Router *mux.Router
}

func (s *Server) Initialize(dbDriver, dbUser, dbPassword, dbPort, dbHost, dbName string) {
	var err error

	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	s.DB, err = gorm.Open(dbDriver, dbURL)
	if err != nil {
		fmt.Printf("cannot connect to %s database", dbDriver)
		log.Fatal("error: ", err)
	}  else {
		fmt.Printf("connected to database..")
	}

	s.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) // database migration
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *Server) Run(addr string) {
	fmt.Println("listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, s.Router))
}