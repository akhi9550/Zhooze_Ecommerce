package main

import (
	"Zhooze/config"
	"Zhooze/db"
	"Zhooze/routes"
	"fmt"
	"log"

	_ "Zhooze/cmd/docs"

	"github.com/gin-gonic/gin"
)

// @title Go + Gin Zhooze E-Commerce API
// @version 1.0.0
// @description Zhooze is an E-commerce platform to purchase and sell shoes
// @contact.name API Support
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:8000
// @BasePath /
// @query.collection.format multi

func main() {
	cfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading the config file")
	}
	fmt.Println(cfig)
	db, err := db.ConnectDatabase(cfig)
	if err != nil {
		log.Fatalf("Error connecting to the database:%v", err)
	}
	router := gin.Default()
	routes.AllRoutes(router, db)

	listenAdder := fmt.Sprintf("%s:%s", cfig.DBPort, cfig.DBHost)
	fmt.Printf("Starting server on %s..\n", cfig.BASE_URL)
	if err := router.Run(cfig.BASE_URL); err != nil {
		log.Fatalf("Error starting server on %s:%v", listenAdder, err)
	}
}
