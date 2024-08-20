package main

import (
	"log"
	"myapp/api/controllers"
	"myapp/api/route"
	"myapp/internal/data"
	"myapp/internal/logic"
	"myapp/internal/session"
	"myapp/internal/users"
	"myapp/pkg/dbconnect"
	"net/http"
)

func main() {
	tarantoolDB, err := dbconnect.ConnectToTarantool()
	if err != nil {
		log.Fatalf("Failed to connect to Tarantool: %v", err)
		return
	}
	defer tarantoolDB.Close()

	sessRepo := session.NewSessionRepo(tarantoolDB)
	dataRepo := data.NewDataRepo(tarantoolDB)
	userRepo := users.NewUserRepo(tarantoolDB)

	logicLayer := logic.NewLogicLayer(dataRepo, sessRepo, userRepo)
	handlerLayer := controllers.NewHandlerLayer(logicLayer)

	router := route.CreateNewRoute(handlerLayer)

	log.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
		return
	}

}
