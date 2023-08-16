package main

import (
	"fmt"
	_db "golang/backend/db"
	_handlers "golang/backend/handlers"
	repo "golang/backend/repositories"
	"golang/backend/services"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	fmt.Println("Server is being started...")
	db := _db.NewDBI()
	err := db.InitDB()
	if err != nil {
		panic(err)
	}
	rateRepository := repo.NewRepository(db.GetDatabase())
	rateService := services.NewService(rateRepository)
	ratesHandler := _handlers.NewHandler(rateService)
	ratesHandler.SyncData()


	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "application/json"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "PUT", "PATCH"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/api/rate/latest", ratesHandler.GetRatesLatest).Methods("GET")
	router.HandleFunc("/api/rate/{date}", ratesHandler.GetRatesByDate).Methods("GET")
	router.HandleFunc("/api/rate/analyze", ratesHandler.GetAnalysisRates).Methods("GET")
	fmt.Println("Server listening on port 8080.")
	_err := http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router))
	if _err != nil {
		panic("Router error: " + err.Error())
	}
	fmt.Scanln()

}
