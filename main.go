package main

import (
	"anonymous/postgres"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)
func main (){ 
	godotenv.Load()
		r := chi.NewRouter()
		port := os.Getenv("PORT")
		databse_url := os.Getenv("DB_URL")
		postgresPool :=postgres.GetconnectionPool(databse_url)
		rows, err := postgresPool.Queryx("SELECT NOW()")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				var currentTime time.Time
				if err := rows.Scan(&currentTime); err != nil {
					log.Fatal(err)
				} 
				log.Println("Current time from database:", currentTime)
			}
	server := http.Server{
			Addr:         net.JoinHostPort("0.0.0.0", port),
			Handler:      r,
			ReadTimeout:  time.Second * 10, 
			WriteTimeout: time.Second * 10,
		}
		log.Println("Server is running on port:", port)
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
}
