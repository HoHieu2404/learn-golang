package cli

import (
	"fmt"
	"os"
    "log"
	_db "learn-golang/back-end/db"
	_handlers "learn-golang/back-end/handlers"
	repo "learn-golang/back-end/repositories"
	"learn-golang/back-end/services"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

    "github.com/urfave/cli/v2"
)

func Run() {
     app := &cli.App{
        Name:  "rateApp",
        Usage: "CLI app for Rate application",
        Commands: []*cli.Command{
            {
                Name:  "start_service",
                Usage: "Run api app",
                Action: func(c *cli.Context) error {
                    StartService()
                    return nil
                },
            },

            {
                Name:  "sync_rate",
                Usage: "Run sync rate data to DB",
                Action: func(c *cli.Context) error {
                    SyncRate()
                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func StartService(){
	db := _db.NewDBI()
	err := db.InitDB()
	if err != nil {
		panic(err)
	}
	rateRepository := repo.NewRepository(db.GetDatabase())
	rateService := services.NewService(rateRepository)
	ratesHandler := _handlers.NewHandler(rateService)

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
}

func SyncRate(){
    db := _db.NewDBI()
	err := db.InitDB()
	if err != nil {
		panic(err)
	}
	rateRepository := repo.NewRepository(db.GetDatabase())
	rateService := services.NewService(rateRepository)
	ratesHandler := _handlers.NewHandler(rateService)
	ratesHandler.SyncData()
}